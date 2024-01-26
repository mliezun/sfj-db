# sfj-db

Single File JSON Database. Store your data in a single JSON file and access it with go structs.

Useful for when you have a very small set of data that you need to read or update.

This projects helps keeping things simple, a plain JSON file that you can read with a text editor.

## Example

Let's say you have a few customers ~50, and are keeping track of all of them in a JSON file.

You have renamed your `pro` plan to `premium` and want to reflect that in your DB.

```go
import "github.com/mliezun/sfj-db"

type Organization struct {
    Id      string `json:"id"`
	Host    string `json:"host"`
	Plan    string `json:"plan"`
}


func UpgradeOrgs() error {
    orgdb, err := sfjdb.Open[[]Organization]("./organizations.json")
    if err != nil {
        return nil, err
    }

    organizations := orgdb.View()
    for _, org := range organizations {
        if org.plan == "pro" {
            org.plan = "premium"
        }
    }

    return orgdb.Save(organizations)
}
```

## Usage

Open a json file as database that is represented by some go type `T`.

```go
db := sfjdb.Open[T](filepath string)
```

Get a copy of the data to do some manipulation.

```go
var mydata T
mydata := db.View()
```

Store a new version of the data. This operation is atomic, it either succeeds (writes successfully to the file) or it fails (it doesn't modify the file at all).

```go
err := db.Save(mydata)
```

## LICENSE

[MIT](/LICENSE)
