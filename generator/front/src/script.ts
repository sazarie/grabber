import { GenerateRandomNumbersRequest } from "./GenerateRandomNumbersRequest";
const socket: WebSocket = new WebSocket("ws://localhost:8080/ws");

const resultElement: HTMLElement = document.getElementById("result")!;

socket.onmessage = function (event: MessageEvent) {
  resultElement.innerHTML += event.data + " ";
};

// sendForm - Функция отправляет форму, получает сгенерированные числа и отображает их на странице
export function sendForm(event: SubmitEvent) {
  event.preventDefault();
  resultElement.innerHTML = "";
  const target: HTMLFormElement = event.target as HTMLFormElement;
  const formData: FormData = new FormData(target);

  const limit: number = parseInt(formData.get("limit")!.toString());
  const goNum: number = parseInt(formData.get("goNum")!.toString());
  console.log(limit);

  const requestData: GenerateRandomNumbersRequest = { limit, goNum };
  console.log(requestData);
  const jsonRequest: string = JSON.stringify({ limit, goNum }); // Создаю строку вида { "limit": 1, "goNum": 1 }
  console.log(jsonRequest);
  socket.send(jsonRequest);
}
