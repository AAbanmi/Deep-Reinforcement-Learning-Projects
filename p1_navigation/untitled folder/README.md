# nhic

NHIC project. In other words, it's Yakeen of Health


## How to run the server

```
$  go run . -config ../etc/stg-config.json
```
You can set the env var in config in your device or pass them like this
```

Prod:
$ STG_DB_HOST=10.10.50.230 STG_DB_USER=webservice_user STG_DB_PASSWORD=r*cNA*G7bbYx STG_DB_NAME=Registry_v2 AUTH_USERNAME=admin AUTH_PASSWORD=admin CALLER_ID=1018851947 STG_GATEWAY_TOKEN=bTAnnO70aApLGr3HiKfLQC6M1SGdMuMW NIC_KEY=bTAnnO70aApLGr3HiKfLQC6M1SGdMuMW NIC_SECRET=UNmv3ni6BBGfW3bB  STG_GATEWAY_URL=https://internal-api.lean.sa   go run . --config  ../etc/stg-config.json 


Testing:
$ STG_DB_HOST=10.12.50.31 STG_DB_USER=Registry_webservice STG_DB_PASSWORD='DtLxD4exZBXQrZBZ' STG_DB_NAME=Registry_v2_dev AUTH_USERNAME=admin AUTH_PASSWORD=admin CALLER_ID=1018851947 STG_GATEWAY_TOKEN=bTAnnO70aApLGr3HiKfLQC6M1SGdMuMW NIC_KEY=bTAnnO70aApLGr3HiKfLQC6M1SGdMuMW NIC_SECRET=UNmv3ni6BBGfW3bB  STG_GATEWAY_URL=https://internal-api.lean.sa   go run . --config  ../etc/stg-config.json 

this should be taken from Saud
caller_id = 1018851947

Server running on port :8000
```

## API Docs
All the APIs document can be found here  `http://localhost:8080/swagger/index.html` 

## Code Documentation



#### Code Structure

``` go 
README.md 
Compat.go  // have compatible citizen structure
config // handles app config
e2e // no use for this pkg as far as i kno
echo  // no use for this pkg as far as i kno
etc // app config
go.mod // app modules
go.sum // app modules 
ihe // no use as far as i know
Nhic.go // handles the business logic here to keep it away from implementation details like http routes
nhic_test.go
Nic // yakeen through MOH only used for Covid19 projects and it is one factor querying where only id is needed.
oauth //  handles getting token from apigee and update it in background
refresh.yml
scfhs // handles the integration with SCFHS
server // pkg when main code(entry point for app) and https routes
store //  handles the logic of talking to our underlying database in this case it's MSSQL Server
yakeen // Yakeen SHC is for registry use only. Lean systems are not allowed to use it and it is NIC direct.

```


#### How to add new endpoints
	1. Go to server/main.go
	2. Add new endpoint as shown below
``` go
// r.Get --> the http method 
// "/organization/{id}" --> endpoint name
// func(w http.ResponseWriter, r *http.Request) --> endpoint method
r.Get("/organization/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        out, err := ctl.GetEstablishment(id)
        if err != nil {
            writeErr(w, err, http.StatusBadRequest)
            return
        }
        writeResponse(w, nil, out)
    })

```


#### Server Errors Standards
If there an error the http status code should be `400` (`http.StatusBadRequest)
``` go 
if err != nil {
                writeErr(w, err, http.StatusBadRequest)
                return
}
```

`writeResponse(w, nil, out)`  will send the  rest with http 200 code 

#### How to add new Swagger doc

1. edit the file under /server/swagger.go
2. run `swag init` to generate json file

#### 3   Party APIs
Stg and prd credentials can be found under Devportal account (email:`NIC Registry@lean.sa`)
and for the development, each developer has to request access for his account on Devportal.

#### How to Run the test 
You can run all test by 
```go 
$ go test -timeout 30s gitlab.lean/leandevclan/nhic/yakeen 
```
s   
or run one test using the method name:
```go 
$ go test -timeout 30s gitlab.lean/leandevclan/nhic/yakeen  -run ^TestYakeen
```


##### Feature Flag
Feature flag to enable and disable ,list of features is enabled in env:
```go 

    if c.FeatureIsEnabled("disable-yakeen") {
        return nil, store.ErrNotFound
    }



// return true if feature in features list
func (c *Controller) FeatureIsEnabled(feature string) bool {
    sort.Strings(c.features)
    i := sort.SearchStrings(c.features, feature)
    if i < len(c.features) && c.features[i] == feature {
        return true
    }
    return false
}




```

#### getFullInfo Endpoint
Once the API is called, it’ll fetch the data in parallel from **getinfo** and **get Contact Info** APIs, then it’ll merge the result and return it.

`waitGroup`  is used to wait for all the goroutines launched here to finish.
 
``` go

    var wg sync.WaitGroup
    var mobileNumber string
    var p *PersonInfo
    var err1 error
    var err2 error

    wg.Add(2)
    go func() {
        mobileNumber, err1 = n.getPersonContactsInfo(id)
        wg.Done()

    }()
    go func() {
        p, err2 = n.getPersonInfo(id)
        wg.Done()

    }()
    wg.Wait()

``` 


#### oauth Pkg 
Oauth pkg deals with Apigee token to get token , saved them locally using bolt database which is  key/value store.
The token are updated in background 

```go 
// any access to access_token and updates are governed by this goroutine
// making sure only one process can control the Token structure
func (o *Oauth) worker() {
    for {
        select {
        case <-time.After(5 * time.Minute):
            for _, i := range *o.consumer {
                tok, err := o.updateToken(i)
                if err != nil && err != ErrNotExpired {
                    log.Println("worker: ", err)
                    continue
                }
                log.Println("Token Updated", tok, err)
            }
        }
    }
}



```


#### Config.json explained

```
{
    "db": {
        "host": "${STG_DB_HOST}",
        "user": "${STG_DB_USER}",
        "password": "${STG_DB_PASSWORD}",
        "port": "${STG_DB_PORT}",
        "name": "${STG_DB_NAME}"
    },
    "gateway": {
        "url": "${STG_GATEWAY_URL}", //APIGEE URL 
        "token": "${STG_GATEWAY_TOKEN}" //APIGEE TOKEN which used in Yakeen
    },
    "oauth": {
        "consumer": [
            {
                "key": "${NIC_KEY}", //New Yakeen 
                "secret": "${NIC_SECRET}",//New Yakeen 
                "name": "nic"
            }
        ],
           "db_path": "./metanode.db"
    },
    "nic": {
        "caller_id": "${CALLER_ID}" //CALLER_ID predefined id
    },
    "auth": { //Basic Auth username and password
        "username": "${AUTH_USERNAME}", // 
        "password": "${AUTH_PASSWORD}"
    },
    "scfhs": {
        "url": "https://internal-api.lean.sa",
        "token": "${STG_SCFHS_TOKEN}" //SCFHS token
    },
    "features": [
        "disable-yakeen",
        "disable-scfhs"
    ]
}


```

### TODO
- Add Prometheus metrics
- Add end-to-end tests in `e2e` package
- Perform load testing


# DELETE record from REG_SMALL
```
-- ****************************************************
-- REMOVE RECORD FROM DB ******************************
-- ****************************************************
-- DELETE FROM REG_SMALL.Individual.Individuals WHERE Health_ID = 'ID10000084583721';
-- SELECT * FROM REG_SMALL.Individual.LuhnNumbersReserve WHERE LuhnNumber = 'ID10000084583721';
-- UPDATE REG_SMALL.Individual.LuhnNumbersReserve SET UsedForIdNumber = null, IsUsed = 0 WHERE LuhnNumber = 'ID10000084583721';
-- SELECT * FROM REG_SMALL.Individual.HealthIDs_NationalIDs_reference WHERE HealthID = 'ID10000084583721';
-- DELETE FROM REG_SMALL.Individual.HealthIDs_NationalIDs_reference WHERE HealthID = 'ID10000084583721';
```
