package plantuml

type RelationType string

const (
	Plain          RelationType = "-"
	Extension      RelationType = "--|>"
	Implementation RelationType = "..|>"
)

type Cardinality string

const (
	ExactlyOne Cardinality = "1"
	ZeroOrOne  Cardinality = "0..1"
	ZeroOrMore Cardinality = "*"
	OneOrMore  Cardinality = "1..*"
)

type Connection struct {
	Source, Target                       string
	SourceCardinality, TargetCardinality Cardinality
	Relation                             RelationType
}

func (c *Connection) Render(t RenderTarget) {
	var relation = "-"
	if c.Relation != "" {
		relation = string(c.Relation)
	}
	var sourceCardinality = ""
	if c.SourceCardinality != "" {
		sourceCardinality = string(c.SourceCardinality)
	}
	var targetCardinality = ""
	if c.TargetCardinality != "" {
		targetCardinality = string(c.TargetCardinality)
	}
	t.P(c.Source, " ", sourceCardinality, relation, targetCardinality, " ", c.Target)
}
