package technical_debt

import (
	"fmt"
	"strings"
)

// graphNodes is the information on a graph.
type graphNode struct {
	name     string
	incoming map[string]bool
	outgoing map[string]bool
}

// String dumps the node into a readable format.
func (n graphNode) String() string {

	var incomingNodeNames []string
	for nodeName := range n.incoming {
		incomingNodeNames = append(incomingNodeNames, nodeName)
	}

	var outgoingNodeNames []string
	for nodeName := range n.outgoing {
		outgoingNodeNames = append(outgoingNodeNames, nodeName)
	}

	return fmt.Sprintf("\nNODE: %s\n\t--> %s\n\t* %s\n\t--> %s \n\n", n.name, strings.Join(incomingNodeNames, ", "), n.name, strings.Join(outgoingNodeNames, ", "))
}

// CalculateDeepCodeFiles dives the codeFiles deep into the stucture.
func CalculateDeepCodeFiles(codeFiles map[string]CodeFile) {

	// Prepare for making a topilogical sort.

	// Create a graph without edges.
	graph := map[string]graphNode{}
	for _, codeFile := range codeFiles {
		node := graphNode{
			name:     codeFile.Name,
			incoming: map[string]bool{},
			outgoing: map[string]bool{},
		}
		graph[node.name] = node
	}

	// Add the edges.
	for _, codeFile := range codeFiles {
		nodeName := codeFile.Name
		node := graph[nodeName]

		// For each edge wire it into the graph.
		for incomingNodeName := range codeFile.DependsOn {

			// Remember this.
			node.incoming[incomingNodeName] = true

			// Complete the reference.
			otherNode := graph[incomingNodeName]
			otherNode.outgoing[nodeName] = true
			graph[incomingNodeName] = otherNode
		}

		// Save our node.
		graph[nodeName] = node
	}

	// Move all nodes that have no incoming edges to a new set.
	noIncomingEdge := map[string]graphNode{}
	for nodeName, node := range graph {
		if len(node.incoming) == 0 {
			noIncomingEdge[nodeName] = node
			delete(graph, nodeName)
		}
	}

	// The total of incoming and no incoming should be the number of inputs to this function.
	if len(noIncomingEdge)+len(graph) != len(codeFiles) {
		panic(`ASSERT: incoming egdge calculations incorrect`)
	}

	// fmt.Printf("\n\n=========== GRAPH\n\n")
	// for _, node := range graph {
	// 	fmt.Println(node)
	// }
	//
	// fmt.Printf("\n\n=========== INCOMING\n\n")
	// for _, node := range noIncomingEdge {
	// 	fmt.Println(node)
	// }

	// Do the topological sort.
	var sortedNodes []graphNode
	for len(noIncomingEdge) > 0 {

		// Grab any single node without incoming edge
		var node graphNode
		for nodeName, nextNode := range noIncomingEdge {
			node = nextNode
			delete(noIncomingEdge, nodeName)
			break // Only want one.
		}

		// Add it to the sorted list.
		sortedNodes = append(sortedNodes, node)

		// Remove this node from each connected node's incoming edge.
		for outgoingNodeName := range node.outgoing {

			// The outgoing node may already have been
			outgoingNode := graph[outgoingNodeName]
			delete(outgoingNode.incoming, node.name)

			// Are there no more incoming nodes?
			if len(outgoingNode.incoming) == 0 {

				// Move this node to the no incoming edge set.
				noIncomingEdge[outgoingNodeName] = outgoingNode
				delete(graph, outgoingNodeName)

			} else {

				// If there are still incoming nodes, put the altered node back.
				graph[outgoingNodeName] = outgoingNode
			}
		}
	}

	// If there is anything left in graph, that is because there is a cycle of dependencies.
	// This is totally fine, in fact it's what we're looking for in this technical debt tool.
	// First use the sorted nodes to fill out dependencies (for efficiency reasons), then
	// iterate through the remaining files in the graph.
	for _, node := range sortedNodes {

		// Get the code file.
		codeFile := codeFiles[node.name]

		// For every outgoing node name, pass on the dependencies.
		for outgoingNodeName := range node.outgoing {
			higherNode := codeFiles[outgoingNodeName]
			for dependsOnFilename := range codeFile.DependsOn {
				higherNode.DependsOn[dependsOnFilename] = true
			}
			codeFiles[outgoingNodeName] = higherNode
		}
	}

	// Now we want to go through the graph (which if there is anything has cycles).
	runAnotherPass := true
	for runAnotherPass {
		runAnotherPass = false
		for _, node := range graph {

			// Get the codeFile node.
			codeFile := codeFiles[node.name]

			// For every outgoing node name, pass on the dependencies.
			for outgoingNodeName := range node.outgoing {
				higherNode := codeFiles[outgoingNodeName]
				for dependsOnFilename := range codeFile.DependsOn {
					// Is this codeFile new?
					if !higherNode.DependsOn[dependsOnFilename] {
						// We are propgating something, keep doing work since it may propogate further.
						runAnotherPass = true
						//fmt.Println("adding " + dependsOnFilename + " to " + node.name)
					}
					higherNode.DependsOn[dependsOnFilename] = true
				}
				codeFiles[outgoingNodeName] = higherNode
			}
		}
	}

	// Ensure that every code file includes what other files depend on it.
	for _, codeFile := range codeFiles {
		codeFile := codeFiles[codeFile.Name]
		for dependsOnFilename := range codeFile.DependsOn {
			// Not that the reverse relation also exists.
			dependsOnNode := codeFiles[dependsOnFilename]
			dependsOnNode.DependedOnBy[codeFile.Name] = true
			codeFiles[dependsOnFilename] = dependsOnNode
		}
	}

	// Lastly, for this technique, every file depends on itself.
	for _, codeFile := range codeFiles {
		codeFile := codeFiles[codeFile.Name]
		codeFile.DependsOn[codeFile.Name] = true
		codeFile.DependedOnBy[codeFile.Name] = true
		codeFiles[codeFile.Name] = codeFile
	}

	// fmt.Println("===========")
	// for _, codeFile := range codeFiles {
	// 	fmt.Printf(codeFile.Name + "\n")
	// 	for filename := range codeFile.DependsOn {
	// 		fmt.Printf("\t" + filename + "\n")
	// 	}
	// 	fmt.Printf("depended on by\n")
	// 	for filename := range codeFile.DependedOnBy {
	// 		fmt.Printf("\t" + filename + "\n")
	// 	}
	// }
}
