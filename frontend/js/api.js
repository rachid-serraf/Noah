
export class NoahAPI {
  constructor(baseURL = "") {
    this.baseURL = baseURL;
  }

  async sendMessage(message) {
    const response = await fetch(`${this.baseURL}/api/chat`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ message }),
    });

    if (!response.ok) {
      throw new Error(`Server returned HTTP ${response.status}`);
    }

    return response.json();
  }
}
