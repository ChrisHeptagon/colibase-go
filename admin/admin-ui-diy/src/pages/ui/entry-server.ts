import { HTMLElement, Window } from "happy-dom";
import { MyElement } from "@/components/testel";

export async function render() {
  const window = new Window();
  const document = window.document;
  // @ts-ignore
  window.customElements.define("my-element", MyElement);
  const myElement = document.createElement("my-element");
  myElement.innerHTML = "Hello World";
  document.body.appendChild(myElement);
  const html = window.document.documentElement.getInnerHTML({includeShadowRoots: true})
  window.happyDOM.cancelAsync();
  return { html };
}
