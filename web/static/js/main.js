function toggleExpandCode() {
  const codeElements = document.getElementsByClassName("code-block");
  const button = document.getElementById("expand-button");
  Array.from(codeElements).forEach((element) => {
    element.classList.toggle("expand");
  });

  let buttonText = "";
  Array.from(codeElements).forEach((element) => {
    if (element.classList.contains("expand")) {
      buttonText = "Collapse Code";
    }
  });
  if (buttonText === "") {
    buttonText = "Expand Code";
  }

  button.innerHTML = buttonText;
}

function toggleDarkMode() {
  const body = document.body;
  body.classList.toggle("dark-mode");

  const button = document.getElementById("dark-mode-button");
  if (body.classList.contains("dark-mode")) {
    button.innerHTML = "Light Mode";
  } else {
    button.innerHTML = "Dark Mode";
  }
}
