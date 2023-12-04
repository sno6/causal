'use client';

import React, { useEffect, useState } from 'react';

import ReactFlow, {
    ConnectionLineType,
    Node,
    Edge,
} from 'reactflow';

import dagre from 'dagre';
import 'reactflow/dist/style.css';

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 172;
const nodeHeight = 36;

const getLayoutedElements = (nodes: any[], edges: any[], direction = 'TB') => {
    dagreGraph.setGraph({ rankdir: direction });

    nodes.forEach((node) => {
        dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
    });

    edges.forEach((edge) => {
        dagreGraph.setEdge(edge.source, edge.target);
    });

    dagre.layout(dagreGraph);

    nodes.forEach((node) => {
        const nodeWithPosition = dagreGraph.node(node.id);
        node.targetPosition = 'top';
        node.sourcePosition = 'bottom';

        // We are shifting the dagre node position (anchor=center center) to the top left
        // so it matches the React Flow node anchor point (top left).
        node.position = {
            x: nodeWithPosition.x - nodeWidth / 2,
            y: nodeWithPosition.y - nodeHeight / 2,
        };

        node.style = { 'border-radius': '100%', 'width': 60, 'height': 60, 'text-align': 'center' };

        return node;
    });

    return { nodes, edges };
};

interface DebugProps {
    clientId: number;
}

export default function Debug({ clientId }: DebugProps) {
    const [initialNodes, setInitialNodes] = useState<Node[]>([]);
    const [initialEdges, setInitialEdges] = useState<Edge[]>([]);

    useEffect(() => {
        setTimeout(() => {
            let inNodes: Node[] = [
                {
                    id: `nid-{"Timestamp":0,"EntityID":${clientId}}`,
                    data: {
                        label: "ROOT"
                    },
                    position: {
                        x: 0,
                        y: 0
                    }
                }
            ];

            let inEdges: Edge[] = [];

            const orderedNodes = window.getNodes(clientId);
            for (let obj of (JSON.parse(orderedNodes) ?? [])) {
                if (obj.removed) {
                    inNodes.push({
                        id: "nid-" + JSON.stringify(obj.id),
                        data: { label: `REMOVED` },
                        position: { x: 0, y: 0 },
                    });
                } else {
                    inNodes.push({
                        id: "nid-" + JSON.stringify(obj.id),
                        data: { label: `node ${obj.value}` },
                        position: { x: 0, y: 0 },
                    });
                }

                inEdges.push({
                    id: "eid-" + JSON.stringify(obj.id),
                    source: "nid-" + JSON.stringify(obj.parent_id),
                    target: "nid-" + JSON.stringify(obj.id),
                    type: ConnectionLineType.Straight,
                    animated: true,
                });
            }

            setInitialNodes(inNodes);
            setInitialEdges(inEdges);
        }, 100);
    })


    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
        initialNodes,
        initialEdges
    );

    return (
        <main className='h-full w-full border-solid border-2 border-ccc-500 mt-5'>
            <ReactFlow
                nodes={layoutedNodes}
                edges={layoutedEdges}
                connectionLineType={ConnectionLineType.Straight}
                fitView
            />
        </main>
    );
};