'use client';

import { SyntheticEvent, useEffect, useRef, useState } from "react";

interface ClientProps {
    num: number;
}

interface TreeNode {
    id: any;
    removed: boolean;
    value: string;
}

const Client = ({ num }: ClientProps) => {
    const [nodes, setNodes] = useState<TreeNode[]>([]);
    const clientId: number = num as number;
    const textAreaRef = useRef(null);

    useEffect(() => {
        setTimeout(() => {
            refreshNodes()
        }, 50)
    })

    const refreshNodes = () => {
        const results = window.getNodes(clientId);
        setNodes(JSON.parse(results) ?? [])
    }

    const nodesToText = () => {
        let stringValue = ""
        for (let obj of nodes) {
            if (obj.removed) {
                continue;
            }
            stringValue += obj.value
        }
        return stringValue;
    }

    const onType = (e: any) => {
        const position = Math.max(0, getCaretPosition()-1);

        const newVal = e.target.value[position];

        let parentId = null
        if (position-1 >= 0 && nodes.length >= position-1) {
            parentId = nodes[position-1].id;
        }

        if (e.nativeEvent.inputType === "deleteContentBackward") {
            if (nodes.length == 1) {
                window.onRemove(nodes[0].id, clientId);
            } else {
                window.onRemove(nodes[position+1].id, clientId);
            }

            return;
        }

        window.onAdd(parentId, newVal, clientId);
    }

    const getCaretPosition = () => {
        if (textAreaRef.current) {
            return (textAreaRef.current as any).selectionStart;
        }
    };

    return (
        <div className="w-150 h-100 flex flex-col m-5">
            <p className="font-semibold text-stone-600">Client {num}</p>
            <textarea 
                ref={textAreaRef}
                value={nodesToText()}
                onChange={(event) => { onType(event) }}
                className="w-100 border border-solid  border-ccc-500 pl-1">
            </textarea>
        </div>
    );
}

export default function Actions() {
    const [clients, setClients] = useState<any[]>([]);

    const addClient = () => {
        const newClientNumber = clients.length + 1;
        window.addClient(newClientNumber);
        
        setClients(prev => {
            return [...prev, 
                <Client key={`client-${newClientNumber}`} num={newClientNumber} />]
        })
    }

    return (
        <div>
            <div className="flex flex-wrap">
                {clients}
            </div>
            <button onClick={addClient} className="w-50 ml-5">+ Add Client</button>
        </div>
    )
}