# Bhojpur CMS - State Transition

The Transition plugin is a finite [*state machine*](https://en.wikipedia.org/wiki/Finite-state_machine) implementation. It can be used standalone, but it integrates nicely with [Bhojpur ORM](https://github.com/bhojpur/orm) models. When integrated with ORM, it will also store state change logs in the database automatically.

## Usage

### Enable Transition for your struct

Embed `transition.Transition` into your struct, it will enable the state machine feature for the struct:

```go
import "github.com/bhojpur/cms/pkg/transition"

type Order struct {
  ID uint
  transition.Transition
}
```

### Define States and Events

```go
var OrderStateMachine = transition.New(&Order{})

// Define initial state
OrderStateMachine.Initial("draft")

// Define a State
OrderStateMachine.State("checkout")

// Define another State and what to do when entering and exiting that state.
OrderStateMachine.State("paid").Enter(func(order interface{}, tx *orm.DB) error {
  // To get order object use 'order.(*Order)'
  // business logic here
  return
}).Exit(func(order interface{}, tx *orm.DB) error {
  // business logic here
  return
})

// Define more States
OrderStateMachine.State("cancelled")
OrderStateMachine.State("paid_cancelled")


// Define an Event
OrderStateMachine.Event("checkout").To("checkout").From("draft")

// Define another event and what to do before and after performing the transition.
OrderStateMachine.Event("paid").To("paid").From("checkout").Before(func(order interface{}, tx *orm.DB) error {
  // business logic here
  return nil
}).After(func(order interface{}, tx *orm.DB) error {
  // business logic here
  return nil
})

// Different state transitions for one event
cancellEvent := OrderStateMachine.Event("cancel")
cancellEvent.To("cancelled").From("draft", "checkout")
cancellEvent.To("paid_cancelled").From("paid").After(func(order interface{}, tx *orm.DB) error {
  // Refund
}})
```

### Trigger an Event

```go
// func (*StateMachine) Trigger(name string, value Stater, tx *gorm.DB, notes ...string) error
OrderStatemachine.Trigger("paid", &order, db, "charged offline by Bhojpur")
// notes will be used to generate state change logs when works with GORM

// When using without Bhojpur ORM, just pass nil to the db, like
OrderStatemachine.Trigger("cancel", &order, nil)

OrderStatemachine.Trigger("cancel", &order, db)
// order's state will be changed to cancelled if current state is "draft"
// order's state will be changed to paid_cancelled if current state is "paid"
```

### Get/Set State

```go
var order Order

// Get Current State
order.GetState()

// Set State
order.SetState("finished") // this will only update order's state, won't save it into database
```

## State change logs

When working with GORM, `Transition` will store all state change logs in the database. Use `GetStateChangeLogs` to get those logs.

```go
// create the table used to store logs first
db.AutoMigrate(&transition.StateChangeLog{})

// get order's state change logs
var stateChangeLogs = transition.GetStateChangeLogs(&order, db)

// type StateChangeLog struct {
//   From       string  // from state
//   To         string  // to state
//   Note       string  // notes
// }
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).