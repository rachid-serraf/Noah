
export class ChatUI {
  constructor(api) {
    this.api = api;

    this.messages = document.querySelector("#messages");
    this.form = document.querySelector("#chat-form");
    this.input = document.querySelector("#message-input");
    this.sendButton = document.querySelector("#send-button");
    this.error = document.querySelector("#error-message");

    this.form.addEventListener("submit", (event) => {
      event.preventDefault();
      this.handleSubmit();
    });
  }

  addMessage(text, sender) {
    const message = document.createElement("div");

    message.classList.add("message", sender);
    message.textContent = text;

    this.messages.appendChild(message);
    this.messages.scrollTop = this.messages.scrollHeight;
  }

  async handleSubmit() {
    const message = this.input.value.trim();

    if (!message) return;

    this.input.value = "";
    this.error.textContent = "";
    this.addMessage(message, "user");

    this.sendButton.disabled = true;
    this.input.disabled = true;

    try {
      const result = await this.api.sendMessage(message);
      this.addMessage(result.reply, "assistant");
    } catch (error) {
      this.error.textContent =
        "Could not reach Noah's server. Check that the backend is running.";
      console.error(error);
    } finally {
      this.sendButton.disabled = false;
      this.input.disabled = false;
      this.input.focus();
    }
  }
}
