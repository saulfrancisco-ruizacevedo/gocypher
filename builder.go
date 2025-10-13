package gocypher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// --- CORE TYPES AND INTERFACES ---

// RelDirection defines the direction of a Cypher relationship.
type RelDirection string

const (
	DirectionOutgoing RelDirection = "-->"
	DirectionIncoming RelDirection = "<--"
	DirectionNone     RelDirection = "--"
)

// PatternPart is an interface representing any component of a Cypher pattern (i.e., a node or a relationship).
type PatternPart interface {
	// render() string
	render(qb *QueryBuilder) string
}

// --- PATTERN STRUCTURES ---

// NodePattern represents a node in a Cypher pattern, e.g., (u:User {name: $name}).
type NodePattern struct {
	Alias      string
	Label      string
	Properties map[string]interface{}
}

// render converts the NodePattern to its Cypher string representation.
// It correctly handles cases where the label is empty (for referencing existing aliases).
func (n *NodePattern) render(qb *QueryBuilder) string {
	labelStr := ""
	if n.Label != "" {
		labelStr = ":" + n.Label
	}

	propStr := ""
	if len(n.Properties) > 0 {
		var props []string
		for key, val := range n.Properties {
			paramNum := qb.paramCounter
			qb.paramCounter++
			paramName := fmt.Sprintf("p%s_%d", paramSanitizer.ReplaceAllString(key, ""), paramNum)
			qb.queryParams[paramName] = val
			props = append(props, fmt.Sprintf("%s: $%s", key, paramName))
		}
		propStr = fmt.Sprintf(" {%s}", strings.Join(props, ", "))
	}

	return fmt.Sprintf("(%s%s%s)", n.Alias, labelStr, propStr)
}

// RelPattern represents a relationship in a Cypher pattern, e.g., -[r:KNOWS {since: 2023}]->.
type RelPattern struct {
	Alias      string
	Type       string
	Direction  RelDirection
	Properties map[string]interface{}
}

// render converts the RelPattern to its Cypher string representation, including properties.
func (r *RelPattern) render(qb *QueryBuilder) string {
	relTypeStr := ""
	if r.Type != "" {
		relTypeStr = ":" + r.Type
	}

	propStr := ""
	if len(r.Properties) > 0 {
		var props []string
		for key, val := range r.Properties {
			paramNum := qb.paramCounter
			qb.paramCounter++
			paramName := fmt.Sprintf("p%s_%d", paramSanitizer.ReplaceAllString(key, ""), paramNum)
			qb.queryParams[paramName] = val
			props = append(props, fmt.Sprintf("%s: $%s", key, paramName))
		}
		propStr = fmt.Sprintf(" {%s}", strings.Join(props, ", "))
	}

	left, right := "-", "-"
	switch r.Direction {
	case DirectionOutgoing:
		right = "->"
	case DirectionIncoming:
		left = "<-"
	}

	return fmt.Sprintf("%s[%s%s%s]%s", left, r.Alias, relTypeStr, propStr, right)
}

// --- FLUENT HELPER FUNCTIONS ---

// N is a shorthand factory function to create a new NodePattern.
func N(alias, label string) *NodePattern {
	return &NodePattern{Alias: alias, Label: label}
}

// NRef is a shorthand factory function to create a new NodePattern.
// NRef references a node by alias only
func NRef(alias string) *NodePattern {
	return &NodePattern{Alias: alias, Label: ""}
}

// WithProperties adds properties to a NodePattern.
func (n *NodePattern) WithProperties(props map[string]interface{}) *NodePattern {
	n.Properties = props
	return n
}

// R is a shorthand factory function to create a new RelPattern.
func R(alias, relType string) *RelPattern {
	return &RelPattern{Alias: alias, Type: relType, Direction: DirectionNone}
}

// WithProperties adds properties to a RelPattern.
func (r *RelPattern) WithProperties(props map[string]interface{}) *RelPattern {
	r.Properties = props
	return r
}

// To sets the relationship's direction to outgoing (e.g., -->).
func (r *RelPattern) To() *RelPattern {
	r.Direction = DirectionOutgoing
	return r
}

// From sets the relationship's direction to incoming (e.g., <--).
func (r *RelPattern) From() *RelPattern {
	r.Direction = DirectionIncoming
	return r
}

// --- QUERY BUILDER ---

var setParamCounter uint64
var paramSanitizer = regexp.MustCompile(`[^a-zA-Z0-9]`)

// QueryBuilder is the main entry point for constructing Cypher queries.
type QueryBuilder struct {
	matchClauses  []string
	createClauses []string
	mergeClauses  []string
	setClauses    []string
	deleteClauses []string
	returnAliases []string
	queryParams   map[string]interface{}
	err           error
	paramCounter  uint64
}

// NewQueryBuilder creates a new instance of the QueryBuilder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		queryParams:  make(map[string]interface{}),
		paramCounter: 0,
	}
}

// renderPattern is an internal helper that renders a pattern and extracts its parameters.
func (qb *QueryBuilder) renderPattern(parts ...PatternPart) string {
	var pattern strings.Builder
	for _, part := range parts {
		pattern.WriteString(part.render(qb))
	}
	return pattern.String()
}

// Match adds a MATCH clause to the query.
func (qb *QueryBuilder) Match(parts ...PatternPart) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	patternStr := qb.renderPattern(parts...)
	if qb.err != nil {
		return qb
	}
	qb.matchClauses = append(qb.matchClauses, "MATCH "+patternStr)
	return qb
}

// OptionalMatch adds an OPTIONAL MATCH clause to the query.
func (qb *QueryBuilder) OptionalMatch(parts ...PatternPart) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	patternStr := qb.renderPattern(parts...)
	if qb.err != nil {
		return qb
	}
	qb.matchClauses = append(qb.matchClauses, "OPTIONAL MATCH "+patternStr)
	return qb
}

// Create adds a CREATE clause to the query.
func (qb *QueryBuilder) Create(parts ...PatternPart) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	patternStr := qb.renderPattern(parts...)
	if qb.err != nil {
		return qb
	}
	qb.createClauses = append(qb.createClauses, "CREATE "+patternStr)
	return qb
}

// Merge adds a MERGE clause to the query.
func (qb *QueryBuilder) Merge(parts ...PatternPart) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	patternStr := qb.renderPattern(parts...)
	if qb.err != nil {
		return qb
	}
	qb.mergeClauses = append(qb.mergeClauses, "MERGE "+patternStr)
	return qb
}

// Set adds a SET clause to update properties.
func (qb *QueryBuilder) Set(updates map[string]interface{}) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	for prop, val := range updates {
		paramNum := qb.paramCounter
		qb.paramCounter++
		paramName := fmt.Sprintf("set%s_%d", paramSanitizer.ReplaceAllString(prop, "_"), paramNum)
		qb.setClauses = append(qb.setClauses, fmt.Sprintf("%s = $%s", prop, paramName))
		qb.queryParams[paramName] = val
	}
	return qb
}

// Delete adds a DELETE clause to the query.
func (qb *QueryBuilder) Delete(aliases ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.deleteClauses = append(qb.deleteClauses, "DELETE "+strings.Join(aliases, ", "))
	return qb
}

// DetachDelete adds a DETACH DELETE clause to delete nodes and their relationships.
func (qb *QueryBuilder) DetachDelete(aliases ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.deleteClauses = append(qb.deleteClauses, "DETACH DELETE "+strings.Join(aliases, ", "))
	return qb
}

// Where adds a WHERE condition.
func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	// Not implemented for brevity, but would be added here.
	return qb
}

// WithParams adds a map of parameters to the query.
func (qb *QueryBuilder) WithParams(params map[string]interface{}) *QueryBuilder {
	// Not implemented for brevity, but would be added here.
	return qb
}

// Return specifies the aliases to be returned by the query.
func (qb *QueryBuilder) Return(aliases ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.returnAliases = append(qb.returnAliases, aliases...)
	return qb
}

// Build validates and assembles the final query string and the parameter map.
func (qb *QueryBuilder) Build() (string, map[string]interface{}, error) {
	if qb.err != nil {
		return "", nil, qb.err
	}
	if len(qb.matchClauses) == 0 && len(qb.createClauses) == 0 && len(qb.mergeClauses) == 0 {
		return "", nil, errors.New("query must have at least one MATCH, CREATE, or MERGE clause")
	}

	var query strings.Builder

	if len(qb.matchClauses) > 0 {
		query.WriteString(strings.Join(qb.matchClauses, "\n") + "\n")
	}
	if len(qb.mergeClauses) > 0 {
		query.WriteString(strings.Join(qb.mergeClauses, "\n") + "\n")
	}
	if len(qb.createClauses) > 0 {
		query.WriteString(strings.Join(qb.createClauses, "\n") + "\n")
	}
	if len(qb.setClauses) > 0 {
		query.WriteString("SET " + strings.Join(qb.setClauses, ", ") + "\n")
	}
	if len(qb.deleteClauses) > 0 {
		query.WriteString(strings.Join(qb.deleteClauses, "\n") + "\n")
	}
	if len(qb.returnAliases) > 0 {
		query.WriteString("RETURN " + strings.Join(qb.returnAliases, ", "))
	}

	return strings.TrimSpace(query.String()), qb.queryParams, nil
}

func PrintQuery(name string, query string, params map[string]interface{}, err error) {
	fmt.Printf("--- %s ---\n", name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Generated Query:")
		fmt.Println(query)
		fmt.Printf("\nParams: %v\n", params)
	}
	fmt.Println("\n" + strings.Repeat("-", 40) + "\n")
}
