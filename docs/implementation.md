# Provider

Provider base project where common events/commands are defined.

Every provider implementation is autonomous service responsible for 3rd party payment gatway integration 
including processing of callbacks and status query requests. 

## Command queue

Every provider must have its own provider command queue. Router SVC will send to this queue all payments request to be processed by this provider. 

Queue name pattern:  provider_%providerName%_commands

    E.g. https://sqs.eu-west-1.amazonaws.com/353673800285/provider_abtek_commands


## Provider event queue

Every provider must publish events with payment status update to common provider events SNS

    E.g. arn:aws:sns:eu-west-1:353673800285:provider_events


## Provider payment process

Provider payment process flow can be described with following steps:

1. Accept payment command from the provider SQS command queue.
2. Process payment with 3rd party gateway.
3. If payment is accepted by / failed at 3rd party gateway, send accepted / failed event to provider event queue.
4. Start scheduled query to 3rd party for payment status update, on succees/fail of payment send success/failed event to provider event queue.
5. Process status update delivered by 3rd party callback, on succees/fail of payment send success/failed event to provider event queue.   

  
### STEP 1 - Accept payment command
 
[ProcessPayOrder command](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-commands/src/main/protobuf/provider-commands.proto) is command that 
triggers processing of payment via provider 3rd party integration. Additional [payment details](https://gitlab.com/carpentumroot/payment-gateway/app/be/router/-/blob/master/router-common/src/main/protobuf/pay-order.proto) description.

Example of ProcessPayOrder command:

    {
        "timestamp":"1618481438601",
        "processPayOrder":{
            "payOrder":{
                "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
                "payMethod":"VAPAY",
                "payType":"PAYIN",
                "amount":"296000",
                "currency":"IDR",
                "merchantCode":"ECLP21-XXX",
                "merchantName":"SBOBet",
                "aggrPayChannelCode":"ECLP21-AHTUMQZVIX_IDR_VAPAY",
                "payChannelCode":"IDR-ABT-VA-MYB-0011",
                "payChannelVertical":"PAY_VERTICAL_UNKNOWN",
                "expirationInSeconds":7200,
                "idPayOrderFromMerchant":"XXXX-YYYYYY-ZZZZZZ",
                "debitor":{
                    "accountName":"First Last",
                    "bankCode":"IDR_016",
                    "bankName":"Maybank Indonesia",
                    "accountType":"UNSET"
                },
                "provider":{
                    "providerCode":"ABTEK",
                    "providerName":"Abtek",
                    "providerAccountCode":"ABTEK-IDR-ABT-0011"
                },
                "merchantFee":{
                    "calcMethod":"PERCENT",
                    "value":"2.0"
                },
                "platformFee":{
                    "calcMethod":"AMOUNT",
                    "value":"4000.0"
                },
                "segment":"SEGMENT_CODE_UNKNOWN",
                "isTest":false
            },
            "acceptedAt":"2021-04-15T10:10:38.55951Z",
            "interfaceSetting":"{}",
            "returnUrl":"http://api.xxxx.com/api/v1/deposit/redirectToMerchant?transactionId=yyyy"
        }
    } 
    
 
### STEP 2 - Process payment with 3rd party gateway
3rd party gateway process may vary, but Provider SVC should be kept as stateless / DB less service. In case of limitation of 3rd party reference ids, dynamo DB table is 
available for ID conversion with following structure:

    payOrderId[String] - Full internal pay order ID 
    provider[String] - provider identification for separating namespace
    payOrderCode[String] - pay order ID code (without tenant code)
    tenantCode[String] - pay order tenant code
    convertedId[String] - converted payorder ID
    ttl[Number] - time to live this item - see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/TTL.html
    
If pay order is accepted by 3rd party, [PayOrderAcceptedByProvider](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event must be published to [provider event queue](#provider-event-queue).

Example of PayOrderAcceptedByProvider event:

    {
        "timestamp":"1618495894009",
        "payOrderAcceptedByProvider":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "acceptedAt":"2021-04-15T14:11:34.009006Z",
            "providerResponse":"VirtualAccountResponse(Some(200),Some(Success))",
            "creditor":{
                "accountNumber":"999111999111",
                "bankCode":"IDR_016",
                "bankName":"Maybank Indonesia",
                "accountType":"UNSET"
            },
            "debitor":{
                "accountName":"Candra Candra",
                "bankCode":"IDR_016",
                "bankName":"Maybank Indonesia",
                "accountType":"UNSET"
            },
            "interfaceSetting":"{}"
        }
    }

In case of any failure  [PayOrderFailed](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event is published to [provider event queue](#provider-event-queue).

Example of published failure event:

    {
        "timestamp":"1618494921176",
        "payOrderFailed":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "failedAt":"2021-04-15T13:55:21.176256Z",
            "failureSource":"GATEWAY_CALLBACK",
            "failureCode":"SYSTEM",
            "failedMessage":"Payment failed via callback",
            "responseCode":"EXPIRED"
        }
    }



### STEP 4 - Provider query scheduling 
For querying provider use provider command queue and event with delayed visibility. As a producer of this event is also a consumer, event structure can be
completely custom, but it is recommended to use common structure defined in [QueryPayOrder command](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-commands/src/main/protobuf/provider-commands.proto)

Result of query may on of following:

1. 3rd party returned success status of transaction, no more queries are scheduled, [PayOrderSucceeded](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event is published to [provider event queue](#provider-event-queue). Origin/authenticity of query response must be resolved. As provider should be stateless, DB less service, 
verifying of payorder amount is left on Payment processor SVC and amount confirmed in query is sent via success event attribute ```forMoney```.

Example of published success event:

    {
        "timestamp":"1618481522752",
        "payOrderSucceeded":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "succeededAt":"2021-04-15T10:12:02.751989Z",
            "forMoney":{"amount":"296000","currency":"IDR"}
        }
    } 

2. 3rd party returned failure status of transaction, no more queries are scheduled, [PayOrderFailed](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event is published to [provider event queue](#provider-event-queue). Origin/authenticity of query response must be resolved.

Example of published failure event:

    {
        "timestamp":"1618494921176",
        "payOrderFailed":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "failedAt":"2021-04-15T13:55:21.176256Z",
            "failureSource":"GATEWAY_CALLBACK",
            "failureCode":"SYSTEM",
            "failedMessage":"Payment failed via callback",
            "responseCode":"EXPIRED"
        }
    }

3. Query must be repeated due to non-terminal status of pay order or connection issues occurred. 
4. Payorder expired and no more queries are scheduled.
      


### STEP 5 - 3rd party callbacks

Recommended URLs for provider callbacks:

    "https://operations.echelonpay.com/provider/%providerName%/payInCallback" 
    "https://operations.echelonpay.com/provider/%providerName%/payOutCallback"


In case of 3rd party callback success status of transaction [PayOrderSucceeded](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event is published to [provider event queue](#provider-event-queue). Origin/authenticity of callback must be resolved. As provider should be stateless, DB less service, 
verifying of payorder amount is left on Payment processor SVC and amount confirmed in query is sent via success event attribute ```forMoney```.

Example of published success event:

    {
        "timestamp":"1618481522752",
        "payOrderSucceeded":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "succeededAt":"2021-04-15T10:12:02.751989Z",
            "forMoney":{"amount":"296000","currency":"IDR"}
        }
    } 

In case of 3rd party callback failed status of transaction [PayOrderFailed](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-events/src/main/protobuf/provider-events.proto) 
event is published to [provider event queue](#provider-event-queue). Origin/authenticity of callback must be resolved.

Example of published failure event:

    {
        "timestamp":"1618494921176",
        "payOrderFailed":{
            "idPayOrder":{"tenantCode":"ECLP","payOrderCode":"ECLP-IDR-D0-111111-2222-3333-4444-5555555"},
            "failedAt":"2021-04-15T13:55:21.176256Z",
            "failureSource":"GATEWAY_CALLBACK",
            "failureCode":"SYSTEM",
            "failedMessage":"Payment failed via callback",
            "responseCode":"EXPIRED"
        }
    }



## Scala implementation
### Provider
Provider should be implementation of trait [Provider](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-service/src/main/scala/com/carpentum/provider/Provider.scala)

### Command/events
For receiving/publishing events class [PayorderReaderSqs](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/blob/master/provider-service/src/main/scala/com/carpentum/provider/payorder/PayorderReaderSqs.scala) is availablee.

### Id conversion
For more convenient way of working with ID conversion table [library](https://gitlab.com/carpentumroot/payment-gateway/lib/be/provider/-/tree/master/provider-persistence-dynamodb/src/main) is available.
Artifact:

    "com.carpentum" %% "provider-persistence-dynamodb" % "1.1.0" 

## Copy provider
To create a new provider, the easiest way is to copy an existing one (**S**) into a new one (**T**).
1. create git repository
1. run `sbt clean` in **S**
1. copy the old repository into a new one `cp -R sbi kotak`
1. run `rename-provider Sbi Kotak` in **T** - script in `utl/be/scripts` to rename folders/files/file content in different uppercase/lowercase variants
1. create the provider using products API `POST /admin/providers`
1. create the interface parameterization (consult with an analyst, see `param.txt` file in products app)
1. create queues in localstack repository to enable local run
1. contact Štěpán to create infrastructure

## Statement processing
Every provider has its own SQS for receiving commands, e.g.

    https://sqs.eu-west-1.amazonaws.com/353673800285/provider_abtek_commands

Information about parsed data is sent into SNS common to all providers, e.g.

    arn:aws:sns:eu-west-1:353673800285:provider_events 

The statement processing can be either manual or automatic.
### Manual processing
1. File is uploaded via API in reconciliation (RC) SVC into S3.
2. Corresponding provider is notified via the SQS.
3. Provider downloads and parses the data.
4. The data is sent into SNS.
5. RC service creates corresponding statements and movements.

### Automatic processing
1. Files are continually uploaded via scraper into an `inbox` folder in S3.
2. RC SVC regularly notifies the provider to check the content of the folder.
3. Provider downloads and parses the files in `inbox` folder and moves them to `archive` folder.
4. The data is sent into SNS.
5. RC service creates corresponding statements and movements.

### Scala implementation
To process the statements in provider:
1. Implement abstract class `FileParser`, mainly the `parse` method to transform file content into the movements. See `StandardFileParser`, `CsvVariantFileParser`, `CsvVariantFileParserWithRepository`.
1. Instantiate `StatementEventsProcessor` class, provide your `FileParser`'s implementation in constructor.