export function printMessage(message: string): void {
  console.log("Message from shared: ", message);
}

export function anotherFunction(): void {
  console.log("Another function in shared package called.");
}
