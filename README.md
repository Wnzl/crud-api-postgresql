### Simple crud api with postgresql and docker containerization

#### Run
```
cp .env.example .env
```
```
docker-compose up -d
```
- Get all users
```
curl --request GET --url localhost:8080/users
```
```
[
   {
      "id":1,
      "first":"John",
      "last":"Doe",
      "created_at":"2020-11-15T01:22:48.221475Z"
   },
   {
      "id":2,
      "first":"James",
      "last":"Doe",
      "created_at":"2020-11-15T01:26:30.210917Z"
   },
   {
      "id":3,
      "first":"Jack",
      "last":"Doe",
      "created_at":"2020-11-15T01:26:36.05569Z"
   }
]
```
- Create new user
```
curl --request POST --url localhost:8080/users --header 'Content-Type: application/json' --data '{"first": "John", "last": "Doe"}'
```
```
{
  "id":1,
  "first":"John",
  "last":"Doe",
  "created_at":"2020-11-15T01:22:48.221475Z"
}
```

- Get specific user
```
curl --request GET --url localhost:8080/users/1
```
```
{
  "id":1,
  "first":"John",
  "last":"Doe",
  "created_at":"2020-11-15T01:22:48.221475Z"
}
```

- Update user
```
curl --request PUT --url localhost:8080/users/1 --header 'Content-Type: application/json' --data '{"first": "Adam", "last": "Wathan"}'
```
```
{
  "id":1,
  "first":"Adam",
  "last":"Wathan",
  "created_at":"2020-11-15T01:22:48.221475Z"
}
```

- Delete user
```
curl --request DELETE --url localhost:8080/users/1
```
```
"successfully deleted"
```
