package main

import "github.com/saulfrancisco-ruizacevedo/gocypher"

func main() {
	// Example 1: Simple - Create a single node
	q1, p1, e1 := gocypher.NewQueryBuilder().
		Create(
			gocypher.N("u", "User").WithProperties(map[string]interface{}{"name": "Alice", "age": 30}),
		).
		Return("u").
		Build()
	gocypher.PrintQuery("Example 1: Create a single node", q1, p1, e1)

	// Example 2: Simple - Find a node
	q2, p2, e2 := gocypher.NewQueryBuilder().
		Match(
			gocypher.N("u", "User").WithProperties(map[string]interface{}{"name": "Alice"}),
		).
		Return("u.name", "u.age").
		Build()
	gocypher.PrintQuery("Example 2: Find a node by property", q2, p2, e2)

	// Example 3: Intermediate - Create a related node and a relationship
	q3, p3, e3 := gocypher.NewQueryBuilder().
		Match(gocypher.N("u", "User").WithProperties(map[string]interface{}{"name": "Alice"})).
		Create(
			gocypher.NRef("u"), // Reference existing 'u'
			gocypher.R("r", "POSTED").To().WithProperties(map[string]interface{}{"date": "2025-09-30"}),
			gocypher.N("p", "Post").WithProperties(map[string]interface{}{"title": "My First Post"}),
		).
		Return("u.name", "r.date", "p.title").
		Build()
	gocypher.PrintQuery("Example 3: Create a relationship", q3, p3, e3)

	// Example 4: Intermediate - Upsert a node and update its properties
	q4, p4, e4 := gocypher.NewQueryBuilder().
		Merge(
			gocypher.N("u", "User").WithProperties(map[string]interface{}{"id": "user123"}),
		).
		Set(map[string]interface{}{
			"u.lastSeen": "2025-09-30T18:00:00Z",
			"u.status":   "active",
		}).
		Return("u").
		Build()
	gocypher.PrintQuery("Example 4: Merge a node and Set properties (Upsert)", q4, p4, e4)

	// Example 5: Complex - Find required and optional data with filtering
	q5, p5, e5 := gocypher.NewQueryBuilder().
		Match(
			gocypher.N("u", "User").WithProperties(map[string]interface{}{"name": "Alice"}),
			gocypher.R("", "POSTED").To(),
			gocypher.N("p", "Post"),
		).
		OptionalMatch(
			gocypher.NRef("p"),
			gocypher.R("", "HAS_COMMENT").From(),
			gocypher.N("c", "Comment"),
		).
		Return("u.name", "p.title", "count(c) AS comments").
		Build()
	gocypher.PrintQuery("Example 5: Complex read with Optional Match", q5, p5, e5)

	// Example 6: Complex - Find and delete a node and its relationships
	q6, p6, e6 := gocypher.NewQueryBuilder().
		Match(
			gocypher.N("p", "Post").WithProperties(map[string]interface{}{"title": "My First Post"}),
		).
		DetachDelete("p").
		Build()
	gocypher.PrintQuery("Example 6: Find and safely delete a node", q6, p6, e6)
}
