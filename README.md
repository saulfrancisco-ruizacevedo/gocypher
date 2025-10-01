# GoCypher

[![Go Reference](https://pkg.go.dev/badge/github.com/saulfrancisco-ruizacevedo/gocypher.svg)](https://pkg.go.dev/github.com/saulfrancisco-ruizacevedo/gocypher)

GoCypher is a lightweight and fluent query builder for [Neo4j](https://neo4j.com/) in Go. It allows you to construct complex Cypher queries programmatically, supporting **nodes, relationships, MATCH, CREATE, MERGE, SET, DELETE, DETACH DELETE**, and more.  

---

## Features

- Fluent and chainable API for building Cypher queries
- Support for nodes and relationships with properties
- Safe parameter handling to prevent injection
- Optional MATCH clauses
- Merge / Upsert nodes
- Set / Update node or relationship properties
- Detach and delete nodes safely

---

## Installation

Import GoCypher in your Go project:

```bash
import "github.com/saulfrancisco-ruizacevedo/gocypher"
```

---

## Quick Example

Here’s a simple example demonstrating how to **create a node** and **return its properties**:

```bash
package main

import (
    "fmt"
    "github.com/saulfrancisco-ruizacevedo/gocypher"
)

func main() {
    // Create a new user node
    query, params, err := gocypher.NewQueryBuilder().
        Create(
            gocypher.N("u", "User")
            .WithProperties(map[string]interface{}{
                "name": "Alice",
                "age":  30,
            }),
        ).
        Return("u").
        Build()

    if err != nil {
        fmt.Println("Error building query:", err)
        return
    }

    fmt.Println("Generated Query:")
    fmt.Println(query)
    fmt.Println("Parameters:", params)
}
```
---

**Output:**

```bash
CREATE (u:User {name: $name, age: $age})
RETURN u
```

**Params:**

```bash
map[string]interface{}{
    "name": "Alice",
    "age":  30,
}
```
---

## Usage Notes

- Use `N(alias, label)` to create new nodes.
- Use `NRef(alias)` to reference an existing node by alias.
- Use `R(alias, type)` to define relationships.
- Use `To()` and `From()` to define relationship direction.
- Chain `Match`, `OptionalMatch`, `Create`, `Merge`, `Set`, `Delete`, and `DetachDelete` to build complex queries.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
