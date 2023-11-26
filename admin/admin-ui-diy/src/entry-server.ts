import { Window } from "happy-dom";
import {testel} from "@/components/testel";

import { MyElement } from "./components/testel";

export async function render() {
  const window = new Window();
  const document = window.document;
  // @ts-ignore
  window.customElements.define("my-element", MyElement);
  document.body.innerHTML = /*html*/`
    <div>
      <h1>Hello World</h1>
      <my-element>
      </my-element>
    </div>
  `;


  const html = window.document.documentElement.getInnerHTML({includeShadowRoots: true})
  return { html };
}
