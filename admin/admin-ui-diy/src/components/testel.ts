
let HTMLElement: typeof window.HTMLElement
let onServer = false
let doc: any
if (typeof window !== "undefined") {
  
  HTMLElement = window.HTMLElement
  console.log("CLIENT ONLY")
  onServer = false
  console.log("onServer", onServer)

  
} else {
  HTMLElement = await import("happy-dom").then(m => m.HTMLElement) as any
  doc = await import("happy-dom").then(m => new m.Window().document)
  console.log("SERVER ONLY")
  onServer = true
  console.log("onServer", onServer)
  
}

export class MyElement extends HTMLElement {
  _count: number = 0
  _html: string = ""
  constructor() {
    super();
    this.attachShadow({ mode: "open" });
  }
  get count() {
    return this._count;
  }
  set count(value) {
    this._count = value;
  }
  get html() {
    return this._html;
  }
  set html(value) {
    this._html = value;
  }
  connectedCallback() {
      if (onServer) {
        console.log("onServer", onServer)
      }
  }
}

