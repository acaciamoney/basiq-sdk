# Basiq.io Golang SDK

This is the documentation for the Golang SDK for Basiq.io API

## Introduction

Basiq.io Golang SDK is a set of tools you can use to easily communicate with Basiq API.
If you want to get familiar with the API docs, [click here](https://basiq.io/api/).

The SDK is organized to mirror the HTTP API's functionality and hierarchy.
The top level object needed for SDKs functionality is the Session
object which requires your API key to be instantiated.
You can grab your API key on the [dashboard](http://dashboard.basiq.io).

## Changelog

1.0.2 - Documentation updated, added RefreshToken method on the Session object

1.0.1 - Documentation updated

1.0.0 - Supported 2.0 API version

0.9.0beta - Initial release

## Getting started

Now that you have your API key, you can use the following command to install the SDK:

```bash
go get -u github.com/basiqio/basiq-sdk-golang/
```

In order to instantiate Session object, following package should be imported:
```go
import (
        "github.com/basiqio/basiq-sdk-golang/basiq"
)
```

Services for API version 1.0 are located in v1 package:
```go
import (
        "github.com/basiqio/basiq-sdk-golang/v1"
)
```

Services for API version 2.0 are located in v2 package:
```go
import (
        "github.com/basiqio/basiq-sdk-golang/v2"
)
```

## Common usage examples

### Fetching a list of institutions

You can fetch a list of supported financial institutions. The function returns a list of Institution structs.

```go
package main

import (
        "github.com/basiqio/basiq-sdk-golang/basiq"
        "log"
)

func main() {
        session, err := basiq.NewSessionV1("YOUR_API_KEY")
        if err != nil {
            log.Printf("%+v", err)
        }

        institutions, err := session.GetInstitutions()
        if err != nil {
            log.Printf("%+v", err)
        }
}
```

Fetching institutions v2.0:

```go
package main

import (
        "github.com/basiqio/basiq-sdk-golang/basiq"
        "log"
)

func main() {
        session, err := basiq.NewSessionV2("YOUR_API_KEY")
        if err != nil {
            log.Printf("%+v", err)
        }

        institutions, err := session.GetInstitutions()
        if err != nil {
            log.Printf("%+v", err)
        }
}
```

### Creating a new connection

When a new connection request is made, the server will create a job that will link user's financial institution with your app. 

```go
package main

import (
        "github.com/basiqio/basiq-sdk-golang/basiq"
        "log"
)

func main() {
        session, err := basiq.NewSessionV1("YOUR_API_KEY")
        if err != nil {
            log.Printf("%+v", err)
        }

        user := session.ForUser(userId)

        job, err := user.CreateConnection(&services.ConnectionData{
                Institution: &services.InstitutionData{
                    Id: "AU00000",
                },
                LoginId:  "gavinBelson",
                Password: "hooli2018",
        })
        if err != nil {
                log.Printf("%+v", err)
        }

        // Poll our server to wait for the credentials step to be evaluated
        connection, err := job.WaitForCredentials(1000, 60)
        if err != nil {
                log.Printf("%+v", err)
        }
}
```

### Fetching and iterating through transactions

In this example, the function returns a transactions list struct which is filtered by the connection.id property. You can iterate 
through transactions list by calling Next().

```go
package main

import (
        "github.com/basiqio/basiq-sdk-golang/basiq"
        "log"
)

func main() {
        session, err := basiq.NewSessionV1("YOUR_API_KEY")
        if err != nil {
            log.Printf("%+v", err)
        }

        user := session.ForUser(userId)

        fb := utilities.FilterBuilder{}
        fb.Eq("connection.id", "conn-id-213-id")
        transactions, err := user.GetTransactions(&fb)
        if err != nil {
                log.Printf("%+v", err)
        }

        for {
                next, err := transactions.Next()
                if err != nil {
                    log.Printf("%+v", err)
                    break
                }

                if next == false {
                    break
                }

                log.Println("Next transactions len:", len(transactions.Data))
        }
}
```

## API

The API of the SDK is manipulated using Services and Entities. Different
services return different entities, but the mapping is not one to one.

### Errors

If an action encounters an error, you will receive an APIError instance.
The struct contains all available data which you can use to act accordingly.

##### APIError struct
```go
type ResponseError struct {
        CorrelationId string `json:"correlationId"`
        Data          []ResponseErrorItem
}

type ResponseErrorItem struct {
        Code   string      `json:"code"`
        Title  string      `json:"title"`
        Detail string      `json:"detail"`
        Source ErrorSource `json:"source"`
}

type APIError struct {
        Data       map[string]interface{}
        Message    string
        Response   ResponseError
        StatusCode int
}
```

Check the [docs](https://basiq.io/api/) for more information about relevant
fields in the error object.

### Filtering

Some of the methods support adding filters to them. The filters are created
using the FilterBuilder struct. After instantiating the struct, you can invoke
methods in the form of Comparison(field, value).

Example:
```go
fb := utilities.FilterBuilder{}
fb.Eq("connection.id", "conn-id-213-id").Gt("transaction.postDate", "2018-01-01")
transactions, err := user.GetTransactions(&fb)
```

This example filter for transactions will match all transactions for the connection
with the id of "conn-id-213-id" and that are newer than "2018-01-01". All you have
to do is pass its reference when you want to use it.


### SDK API List

<details>
<summary>
Services
</summary>

#### Session

##### Creating a new Session object

##### Version 1.0.
```go
session, err := basiq.NewSessionV1("YOUR_API_KEY")
```

##### Version 2.0.
```go
session, err := basiq.NewSessionV2("YOUR_API_KEY")
```

##### Refreshing a token
```go
err := session.RefreshToken()
```

##### Getting institutions
```go
institutions, err := session.GetInstitutions()
```

##### Getting institution
```go
institution, err := session.GetInstitution("INSTITUTION_ID")
```

##### Creating a new user
```go
user, err := session.CreateUser(userData)
```

##### Referencing a user
*Note: The following action will not send an HTTP request, and can be used
to perform additional actions for the instantiated user.*

```go
user := session.ForUser(userId)
```

#### UserService

The following are APIs available for the User service

##### Creating a new UserService

```go
userService := v1.NewUserService(session)
```

##### Version 2.0.
```go
userService := v2.NewUserService(session)
```

##### Referencing a user
*Note: The following action will not send an HTTP request, and can be used
to perform additional actions for the instantiated user.*

```go
user := userService.ForUser(userId)
```

##### Creating a new User

```go
user, err := userService.CreateUser(&Services.UserData{
        Mobile: "+61410888555",
})
```

##### Getting a User

```go
user, err := userService.GetUser(userId)
```

##### Update a User

```go
user, err := userService.UpdateUser(userId, &Services.UserData{})
```

##### Delete a User

```go
err := userService.DeleteUser(userId)
```

##### Refresh connections

```go
err := userService.RefreshAllConnections(userId)
```

##### List all connections

```go
conns, err := userService.ListAllConnections(userId, *filter)
```

##### Get account

```go
acc, err := userService.GetAccount(userId, accountId)
```

##### Get accounts

```go
accs, err := userService.GetAccounts(userId, *filter)
```

##### Get transaction

```go
transaction, err := userService.GetTransaction(userId, transactionId)
```

##### Get transactions

```go
transactions, err := userService.GetTransactions(userId, *filter)
```

#### ConnectionService

The following are APIs available for the Connection service

##### Creating a new ConnectionService

```go
connService := v1.NewConnectionService(session, user)
```

##### Version 2.0.
```go
connService := v2.NewConnectionService(session, user)
```

##### Get connection

```go
connection, err := connService.GetConnection(connectionId)
```

##### Get connection entity with ID without performing an http request

```go
connection := connService.ForConnection(connectionId)
```

##### Create a new connection

```go
job, err := connService.NewConnection(*connectionData)
```

##### Update connection

```go
job, err := connService.UpdateConnection(connectionId, password)
```

##### Delete connection

```go
err := connService.DeleteConnection(connectionId)
```

##### Get a job

```go
job, err := connService.GetJob(jobId)
```


#### TransactionService

The following are APIs available for the Transaction service

##### Creating a new TransactionService

```go
transactionService := v1.NewTransactionService(session)
```

##### Version 2.0.
```go
transactionService := v2.NewTransactionService(session)
```

##### Get transactions

```go
transactionList, err := transactionService.GetTransactions(userId, *filter)
```

#### InstitutionService

The following are APIs available for the Institution service

##### Creating a new InstitutionService

```go
instService := v1.NewInstitutionService(session, userId)
```

##### Version 2.0.
```go
instService := v2.NewInstitutionService(session, userId)
```

##### Get institutions

```go
institutions, err := instService.GetInstitutions()
```

##### Get institution

```go
institution, err := instService.GetInstitution(institutionId)
```

</details>


<details><summary>
Entities
</summary>

##### Updating a user instance [mut]

```go
err := user.Update(&Services.UserData{
        Mobile: "+61410888665",
})
```

##### Deleting a user

```go
err := user.Delete()
```

##### Get all of the user's accounts

```go
accounts, err := user.GetAccounts()
```

##### Get a user's single account

```go
account, err := user.GetAccount(accountId)
```

##### Get all of the user's transactions

```go
transactions, err := user.GetTransactions()
```

##### Get a user's single transaction

```go
transaction, err := user.GetTransaction(transactionId)
```

##### Create a new connection

```go
job, err := user.CreateConnection(&services.ConnectionData{
         Institution: &services.InstitutionData{
             Id: "AU00000",
         },
         LoginId:  "gavinBelson",
         Password: "hooli2018",
})
```

##### Refresh all connections

```go
err := user.RefreshAllConnections()
```

#### Connection

##### Refresh a connection

```go
job, err := connection.Refresh()
```

##### Update a connection

```go
job, err := connection.Update(password)
```

##### Delete a connection

```go
err := connection.Delete()
```

#### Job

##### Get the connection id (if available)

```go
connectionId := job.GetConnectionId()
```

##### Get the connection

```go
connection, err := job.GetConnection()
```

##### Get the connection after waiting for credentials step resolution
(interval is in milliseconds, timeout is in seconds)

```go
connection, err := job.WaitForCredentials(interval, timeout)
```

##### Get the connection after waiting for transactions step resolution
(interval is in milliseconds, timeout is in seconds)

```go
connection, err := job.WaitForTransactions(interval, timeout)
```

#### Transaction list

##### Getting the next set of transactions [mut]

```go
next, err := transactions.Next()
```
</details>