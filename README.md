# Go JSON v2 Playground

A demonstration project showcasing the new `encoding/json/v2` package in Go 1.25, comparing default JSON marshaling with custom marshalers for complex data types.

## Features

- **Default JSON Encoding**: Standard marshaling behavior for `time.Time` and `*big.Float`
- **Custom JSON Encoding**: Custom marshalers for RFC3339 time formatting and big.Float as strings
- **HTTP Server**: Two endpoints demonstrating different encoding approaches

## Requirements

- Go 1.25 or later (required for `encoding/json/v2` package)

## Installation

```bash
git clone <repository-url>
cd go-jsonv2-playground
go mod tidy
go run cmd/main.go
```

## Usage

Start the server:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080` with two available endpoints.

### Endpoints

#### `/default` - Default JSON Encoding

Uses standard JSON marshaling without custom formatters.

**Example Request:**
```bash
curl http://localhost:8080/default
```

**Example Response:**
```json
{
  "id": 1,
  "timestamp": "2024-08-27T14:30:45.123456789Z",
  "value": 12345.6789,
  "message": "default encoding"
}
```

#### `/custom` - Custom JSON Encoding

Uses custom marshalers for enhanced formatting:
- `time.Time` formatted as RFC3339 strings
- `*big.Float` values as quoted strings for precision
- Zero struct fields omitted

**Example Request:**
```bash
curl http://localhost:8080/custom
```

**Example Response:**
```json
{
  "id": 2,
  "timestamp": "2024-08-27T14:30:45Z",
  "value": "98765.4321",
  "message": "custom encoding"
}
```

## Code Structure

### Data Structure

```go
type Payload struct {
    ID        int        `json:"id"`
    Timestamp time.Time  `json:"timestamp"`
    Value     *big.Float `json:"value"`
    Message   string     `json:"message"`
}
```

### Custom Marshalers

The project demonstrates two custom marshalers:

1. **Time RFC3339 Marshaler**: Formats `time.Time` as RFC3339 strings
2. **Big Float String Marshaler**: Converts `*big.Float` to quoted string representation

```go
var timeRFC3339 = json.MarshalFunc(func(t time.Time) ([]byte, error) {
    return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
})

var bigFloatAsString = json.MarshalFunc(func(f *big.Float) ([]byte, error) {
    return []byte(`"` + f.Text('g', -1) + `"`), nil
})
```

## JSON v2 Features Demonstrated

### Default vs Custom Marshaling

- **Default endpoint**: Uses `json.Marshal(p, json.StringifyNumbers(false))`
- **Custom endpoint**: Uses multiple marshaling options:
  ```go
  json.Marshal(p,
      json.WithMarshalers(timeRFC3339),
      json.WithMarshalers(bigFloatAsString),
      json.OmitZeroStructFields(true),
  )
  ```

### Key Differences

| Feature | Default | Custom |
|---------|---------|--------|
| Time Format | Go's default time format | RFC3339 string |
| Big Float | Numeric representation | Quoted string |
| Zero Fields | Included | Omitted |
| Precision | May lose precision | Full precision preserved |

## Testing Examples

### Compare Outputs

You can easily compare the outputs by running both endpoints:

```bash
# Default encoding
curl -s http://localhost:8080/default | jq '.'

# Custom encoding  
curl -s http://localhost:8080/custom | jq '.'
```

### Script for Testing

Create a simple test script:

```bash
#!/bin/bash
echo "=== Default Encoding ==="
curl -s http://localhost:8080/default | jq '.'

echo -e "\n=== Custom Encoding ==="
curl -s http://localhost:8080/custom | jq '.'
```

## Use Cases

This playground is useful for understanding:

- Migration from `encoding/json` to `encoding/json/v2`
- Custom marshaling strategies for complex types
- Precision handling with `big.Float` in JSON
- Time formatting options in JSON APIs
- Performance comparisons between encoding approaches
