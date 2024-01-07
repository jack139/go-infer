
## API documentation (template)

### 1. Global parameters definition

Input parameters

| Parameters | Type   | Explanation            | Example |
| ---------- | ------ | ---------------------- | ------- |
| appId     | string | application ID          |         |
| version   | string | version number          |         |
| signType  | string | Signature algorithm     | SM2 or SHA256 |
| signData  | string | Signature data, see below |            |
| encType   | string | Data encryption algorithm | plain       |
| timestamp | int    | unix timestamp (seconds)  |             |
| data      | json   | Data, see each call definition for details  |         |

> Signature/signature verification algorithm:
>
> 1. Filter: Get the parameter key-value pairs and eliminate the three parameters signData, encData, and extra. The data parameters are arranged in ascending order by key for json serialization.
> 2. Sort: Sort by key in ascending order.
> 3. Splicing: splicing request parameters in sorted order.
>
> ```key1=value1&key2=value2&...&key=appSecret```, key=appSecret is fixed at the end of the parameter string, and appSecret needs to be replaced with the appSecret assigned by the application channel.
>
> 4. Signature: Use the specified algorithm to sign to obtain binary bytes, use hexadecimal encoding to obtain the signature string, and then base64 encode.
> 5. Signature verification: Sign the received parameters according to steps 1-4, and compare whether the obtained signature string is consistent with the submitted signature string.

Signature example:

```
Request parameters:
{
    "appId":"3EA25569454745D01219080B779F021F",
    "version": "1",
    "signType": "SM2",
    "signData": "...",
    "encType": "plain",
    "timestamp":1658716494,
    "data": {
        "text":"test text",
        "image":""
    }
}

Keys:
appSecret="41DF0E6AE27B5282C07EF5124642A352"
SM2_privateKey="JShsBOJL0RgPAoPttEB1hgtPAvCikOl0V1oTOYL7k5U="

String to be signed:
appId=3EA25569454745D01219080B779F021F&data={"image":"","text":"test text"}&encType=plain&signType=SHA256&timestamp=1658716494&version=1&key=41DF0E6AE27B5282C07EF5124642A352

SHA256 signature result:
"a68c1b852a650314afaad684f3652c336c9b969e943825a29380b516de746ece"

The result after base64:
"YTY4YzFiODUyYTY1MDMxNGFmYWFkNjg0ZjM2NTJjMzM2YzliOTY5ZTk0MzgyNWEyOTM4MGI1MTZkZTc0NmVjZQ=="

SM2 signature result (different each time):
"ILSOY5A0/sfW5Y9T6rIjl1AEPlDtQeqtwAxLibNbnajlj2fY/DxvTuSok+sqxy2St4pvvs4/rdaNOCNpwBuJ6A=="
```

Return results

| Parameters | Type   | Explanation            | Example |
| --------- | ------- | ------------------------------------------------------------ | ----- |
| appId     | string  | application ID                                               |       |
| code      | string  | returned status code                                         |       |
| signType  | string  | Signature algorithm (no signature required)                 | plain |
| encType   | string  | Data encryption algorithm (currently not encrypted)         | plain |
| success   | boolean | success or failure                                          |       |
| timestamp | int     | unix timestamp                                        |       |
| requestId | string  | The ID of the current request                         |       |
| data      | json    | data is returned when success, data.msg returns an error message when an error occurs |       |

> When successful: code is 0, success is true, and the data content is defined in each interface;
> When an error occurs: code returns an error code. For specific definitions, see the description of each interface.

Return example

```json
{
    "appId": "3EA25569454745D01219080B779F021F", 
    "code": 0, 
    "signType": "plain",
    "encType": "plain",
    "success": true,
    "timestamp": 1658716495,
    "data": {
       "msg": "success", 
       ...
       "requestId": "20220727e5b9bbfe0e33c01e3c8ccb9c7382d512"
    }
}
```

Global error code

| Error code | Explanation                  |
| ---- | ---------------------------------- |
| 9800 | Invalid signature                   |
| 9801 | An error in the signature parameters  |
| 9802 | Unixtime is outside the acceptable range |



### 2. Get text features in example

> Get text embeddings
>
> Note:
>
> 1. This is just a demonstration, the results are not guaranteed to be accurate.
> 2. Model weights are exported from the official Chinese Bert, see the export directory for details.

Request URL

> http://127.0.0.1:5000/api/embedding

Request method

> POST

Input parameters

| Parameters | Required   | Type | Explanation   |
| ----- | ---- | ------ | ------------------ |
| text | Yes   | string | Input text         |

Request example

```json
{
    "text" : "test text"
}
```

Return results

| Parameters | Required   | Type | Explanation   |
| ---------- | ---- | ------ | -------------------- |
| embeddings     | Yes   | float array | text embeddings |


Return example

```json
{
    "appId": "", 
    "code": 0, 
    "data": {
        "embeddings": [
            0.21856225, 0.22315553, 
            ...
            , -0.5392351, -0.14255117
        ], 
        "msg": "success",
        "requestId": "2022072792c2e34ae170db21066849f015dd3133"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1658716495
}
```

### 3. Get image features in example

> Use Mobilenet to get image embeddings
>
> Note:
>
> 1. This is just a demonstration, the results are not guaranteed to be accurate.
> 2. Model weights are exported from Keras official weights, see the export directory for details.

Request URL

> http://127.0.0.1:5000/api/mobile

Request method

> POST

Input parameters

| Parameters | Required   | Type | Explanation   |
| ----- | ---- | ------ | ------------------ |
| image | Yes   | string | base64 encoded image data |
> size of the picture: 224\*224

Request example

```json
{
    "image" : "..."
}
```

Return results

| Parameters | Required   | Type | Explanation   |
| ---------- | ---- | ------ | -------------------- |
| embeddings     | Yes   | float array | image embeddings |


Return example

```json
{
    "appId": "", 
    "code": 0, 
    "data": {
        "embeddings": [[[
            1.802082e-11, 2.1090724e-10, 
            ..., 
            1.6219401e-10, 1.038099e-10
        ]]], 
        "msg": "success",
        "requestId": "20220727b796f2356b94bcbd47679b2606a8b117"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1658900086
}
```
