package data

import (
	"fmt"
	"math"
	"strings"
)

type Node struct {
	ID    string
	X     float64
	Y     float64
	Vx    float64
	Vy    float64
}

type Edge struct {
	Source *Node
	Target *Node
}

func (gen *Generator) generateKBGraph() string {
	nodes := make(map[string]*Node)
	var edges []Edge

	for slug := range gen.KBs {
		if slug == "front" {
			continue
		}
		nodes[slug] = &Node{
			ID: slug,
			X:  math.Cos(float64(len(nodes))) * 100,
			Y:  math.Sin(float64(len(nodes))) * 100,
		}
	}

	for slug, kb := range gen.KBs {
		if slug == "front" {
			continue
		}
		for _, link := range kb.Backlinks {
			if nodes[link.Slug] != nil && nodes[slug] != nil {
				edges = append(edges, Edge{Source: nodes[link.Slug], Target: nodes[slug]})
			}
		}
	}

	// Simple Fruchterman-Reingold algorithm
	iterations := 100
	area := 800.0 * 600.0
	k := math.Sqrt(area / float64(len(nodes)+1))
	t := 100.0
	dt := t / float64(iterations+1)

	for i := 0; i < iterations; i++ {
		for _, v := range nodes {
			v.Vx = 0
			v.Vy = 0
			for _, u := range nodes {
				if u != v {
					dx := v.X - u.X
					dy := v.Y - u.Y
					dist := math.Sqrt(dx*dx + dy*dy)
					if dist > 0 {
						repel := (k * k) / dist
						v.Vx += (dx / dist) * repel
						v.Vy += (dy / dist) * repel
					}
				}
			}
		}

		for _, e := range edges {
			dx := e.Target.X - e.Source.X
			dy := e.Target.Y - e.Source.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 0 {
				attract := (dist * dist) / k
				e.Source.Vx += (dx / dist) * attract
				e.Source.Vy += (dy / dist) * attract
				e.Target.Vx -= (dx / dist) * attract
				e.Target.Vy -= (dy / dist) * attract
			}
		}

		for _, v := range nodes {
			dist := math.Sqrt(v.Vx*v.Vx + v.Vy*v.Vy)
			if dist > 0 {
				v.X += (v.Vx / dist) * math.Min(dist, t)
				v.Y += (v.Vy / dist) * math.Min(dist, t)
			}
			v.X = math.Max(-400, math.Min(400, v.X))
			v.Y = math.Max(-300, math.Min(300, v.Y))
		}
		t -= dt
	}

	var sb strings.Builder
	sb.WriteString(`<svg viewBox="-450 -350 900 700" xmlns="http://www.w3.org/2000/svg" style="width:100%;height:500px;background:var(--bg);border: 1px solid var(--border);border-radius:8px;margin-bottom:2rem;">`)
	sb.WriteString(`<g stroke="var(--text-light)" stroke-width="1">`)
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf(`<line x1="%f" y1="%f" x2="%f" y2="%f" opacity="0.3"/>`, e.Source.X, e.Source.Y, e.Target.X, e.Target.Y))
	}
	sb.WriteString(`</g>`)
	sb.WriteString(`<g stroke="var(--bg)" stroke-width="1.5">`)
	for _, v := range nodes {
		sb.WriteString(fmt.Sprintf(`<circle cx="%f" cy="%f" r="6" fill="var(--accent)"/>`, v.X, v.Y))
		sb.WriteString(fmt.Sprintf(`<a href="/kb/%s"><text x="%f" y="%f" font-family="sans-serif" font-size="12px" fill="var(--text)" stroke="none" transform="translate(10, 4)">%s</text></a>`, v.ID, v.X, v.Y, v.ID))
	}
	sb.WriteString(`</g></svg>`)

	return sb.String()
}
