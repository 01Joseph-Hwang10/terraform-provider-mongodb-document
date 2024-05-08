data "mongodb_database" "default" {
  name = "default"
}

data "mongodb_database_collection" "users" {
  database = data.mongodb_database.default.name
  name     = "users"
}

data "mongodb_database_index" "user_age_index" {
  database   = data.mongodb_database.default.name
  collection = data.mongodb_database_collection.users.name
  // Assuming that the index is created on the "age" field
  // with an ascending direction, the index name generated by
  // MongoDB will be "age_1"
  index_name = "age_1"
}