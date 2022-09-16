# Maeve: hierarchical configuration made simple

Maeve is a plug-and-play configuration system for Postgres that aims to be correct, fast and simple to reason about.

## Features

- Simple, hierarchical and performant key-value path in Postgres.
- Correct updates and deletes with minimal effort: optimistic locking is used to ensure no value gets overriden. No race conditions!
- Support multiple different use cases: settings per carrier, per user, namespaces
- Go API for reading and writing settings

## Example

```go
// first, create a maeve client:
import (
  "github.com/nubunto/maeve"
  "github.com/nubunto/maeve/pg"
  "github.com/nubunto/maeve/mavemock"
)

mcli, _ := pg.NewClient()

// read all settings under the path "foo/*"
// invoiceSettings is a maeve.KeyValueList
invoiceSettings, err := mcli.Fetch(context.TODO(), maeve.Path("invoices/*"))

// Put creates the values provided in the list in bulk
// number of keys and values must match, otherwise an error is returned
// by default, Put will override values
err := mcli.Put(maeve.KV("invoices/some_property", "value 1", "invoices/other_property", "value 2"))

// you can delete a property as well:
err := mcli.Delete(maeve.Path("invoices/some_property"))

// you can use the maevemock package for unit tests
ctrl := gomock.NewController(t)
maeveMockClient := maevemock.NewMockClient(ctrl)
maeveMockClient.Mock.EXPECT().
  Put(maevemock.KVMatcher("invoices/some_property", "value 1", "invoices/other_property", "value 2"))

err := maeveMockClient.Put(maeve.KV("invoices/some_property", "value 1", "invoices/other_property", "value 2"))
```

## Hooks

Hooks are functions that are fired in certain events for key-value pairs. You can register hooks like this:

```go
mcli.RegisterHook(&maeve.Hook{
  Path: "invoices/*",
  OnCreate: func(props maeve.KeyValueList) {
    for _, prop := range props {
      fmt.Println("created", prop.Key, "with value", prop.Value)
    }
  },
  OnUpdate: func(props maeve.KeyValueList) {
    for _, prop := range props {
      fmt.Println("updated", prop.Key, "with value", prop.Value)
    }
  },
  OnDelete: func(props maeve.KeyValueList) {
    for _, prop := range props {
      fmt.Println("deleted", prop.Key, "with value", prop.Value)
    }
  },
})
```

`OnCreate` is called when the key-value pair did not exist before.

`OnUpdate` is called when the key-value pair existed before, and right after it is updated.

`OnDelete` is called right after the key-value pair is deleted.

When used with a glob path, hooks will be called once with a list of affected key-pairs.

There are no guarantees that the values will still be valid by the time hooks are called. Thus, hooks are supposed to be idempotent operations.

## Correct updates in the face of concurrency

Updates to settings will use optimistic locking. This means that concurrent updates to a given path will be guaranteed to be processed in the order they were received.
