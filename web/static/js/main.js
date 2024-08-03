function toggleExpandCode() {
  const codeElements = document.getElementsByClassName("code-block");
  console.log("Code elements: ", codeElements);
  Array.from(codeElements).forEach((element) => {
    element.classList.toggle("expand");
  });
}

function toggleDarkMode() {
  const body = document.body;
  body.classList.toggle("dark-mode");
}
