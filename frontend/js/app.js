import { NoahAPI } from "./api.js";
import { ChatUI } from "./chat.js";

class NoahApp {
  constructor() {
    this.api = new NoahAPI();
    this.chat = new ChatUI(this.api);
  }
}

new NoahApp();
