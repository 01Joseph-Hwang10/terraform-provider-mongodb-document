---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mongodb_database_document Data Source - mongodb"
subcategory: ""
description: |-
  This resource creates a single document in a collection
  in a database on the MongoDB server.
---

# mongodb_database_document (Data Source)

This resource creates a single document in a collection 
in a database on the MongoDB server.

## Example Usage

```terraform
data "mongodb_database" "default" {
  name = "default"
}

data "mongodb_database_collection" "users" {
  database = data.mongodb_database.default.name
  name     = "users"
}


data "mongodb_database_document" "first_user" {
  database    = data.mongodb_database.default.name
  collection  = data.mongodb_database_collection.users.name
  document_id = "<stringified-mongodb-object-id>"
}

// Example usage of the data source: writing the document to a local file
resource "local_file" "first_user_document" {
  content  = data.mongodb_database_document.first_user.document
  filename = "${path.module}/first-user.json"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `collection` (String) Name of the collection to read the document in.
- `database` (String) Name of the database to read the collection in.
- `document_id` (String) <p>Document ID of the document.</p>  <p>This value is a stringified MongoDB ObjectID.</p>  <p>In golang, you can use the following code to stringify an ObjectID:</p>  <pre><code class="language-go">objectID.(primitive.ObjectID).Hex()</code></pre>

### Read-Only

- `document` (String) <p>Document to insert into the collection.</p>  <p>The value of this attribute is a stringified JSON, with every double quote escaped with a backslash. This means that the JSON string contains backslashes before every double quote.</p>  <p>In terraform, you&rsquo;ll be able to smoothly decode the JSON string by using the <code>jsondecode</code> function.</p>  <pre><code class="language-terraform">decoded = jsondecode(document)</code></pre>
- `id` (String) <p>Resource identifier.</p>  <p>ID has a value with a format of the following:</p>  <pre><code class="">databases/&lt;database&gt;/collections/&lt;name&gt;/documents/&lt;document_id&gt;</code></pre>