# optional
Optionals implementation for Golang

Could help to reduce amount of annoying 
```go
result, err := doSomething()
if err != nil {
  return nil, err
}
nextResult, nextErr := doSomethingElse(result)
if nextErr != nil {
  return nil, nextErr
}
finalResult, finalErr := doSomethingOneMoreTime(nextResult)
if finalErr != nil {
  return nil, finalErr
}
return finalResult
```

to 

```go
return optional.MapErr(
  optional.MapErr(
    optional.OfError(doSomething()),
    doSomethingElse,
  ),
  doSomethingOneMoreTime).Get()
```
