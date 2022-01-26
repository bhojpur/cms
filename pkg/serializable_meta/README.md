# Bhojpur CMS - Serializable Meta

The Serializable Meta allows the developer to specify, for a given Model, a custom serialization model along with field mappings. This mechanism thus allows one model to act as another model when it comes to serialization.

## Usage

The example herein shows how to manage different kinds of background jobs using https://github.com/bhojpur/cms/pkg/serializable_meta.

### Define serializable model

Define `BhojpurJob` model and embed `serializable_meta.SerializableMeta` to apply the feature.

```go
type BhojpurJob struct {
  gorm.Model
  Name string
  serializable_meta.SerializableMeta
}
```

Add function `GetSerializableArgumentResource` to the model, so [Serializable Meta](https://github.com/bhojpur/cms/pkg/serializable_meta) can know the type of argument. Then, define background jobs.

```go
func (bhojpurJob BhojpurJob) GetSerializableArgumentResource() *admin.Resource {
  return jobsArgumentsMap[bhojpurJob.Kind]
}

var jobsArgumentsMap = map[string]*admin.Resource{
  "newsletter": admin.NewResource(&sendNewsletterArgument{}),
  "import_products": admin.NewResource(&importProductArgument{}),
}

type sendNewsletterArgument struct {
  Subject string
  Content string
}

type importProductArgument struct {}
```

### Use serializable features

At first, set a Job's `Name`, `Kind` and `SetSerializableArgumentValue`. Then, save it into database.

```go
var bhojpurJob BhojpurJob
bhojpurJob.Name = "sending newsletter"
bhojpurJob.Kind = "newsletter"
bhojpurJob.SetSerializableArgumentValue(&sendNewsletterArgument{
  Subject: "subject",
  Content: "content",
})

db.Create(&bhojpurJob)
```

This will marshal `sendNewsletterArgument` as a json, and save it into database by this SQL

```sql
INSERT INTO "bhojpur_jobs" (kind, value) VALUES (`newsletter`, `{"Subject":"subject","Content":"content"}`);
```

Now you can fetch the saved `BhojpurJob` from the database. And get the serialized data from the record.

```go
var result BhojpurJob
db.First(&result, "name = ?", "sending newsletter")

var argument = result.GetSerializableArgument(result)
argument.(*sendNewsletterArgument).Subject // "subject"
argument.(*sendNewsletterArgument).Content // "content"
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).