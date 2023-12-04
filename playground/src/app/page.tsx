'use client';

import Actions from "./actions";
import Debug from "./debug";
import { useEffect } from "react";

declare var Go: any

declare global {
  interface Window {
    getNodes: (clientId: number) => string;
    onAdd(parentId: any, newVal: string, clientId: number): void;
    onRemove(nodeId: any, clientId: number): void;
    addClient(clientId: number): void;
  }
}

export default function Home() {
  useEffect(() => {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
      go.run(result.instance)
    });
  }, [])

  return (
    <main className='h-full w-full'>
      <Actions></Actions>
      <Debug clientId={1}></Debug>
    </main>
  );
};