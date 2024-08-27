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
      buttonText = loadSvg("minimize-2", 14);
    }
  });
  if (buttonText === "") {
    buttonText = loadSvg("maximize-2", 14);
  }

  button.innerHTML = buttonText;
}

function loadSvg(name, size) {
  return `<img src="/static/svg/${name}.svg" alt="${name}" class="icon"/>`;
}

function toggleDarkMode() {
  const body = document.body;
  const button = document.getElementById("dark-mode-button");

  // Toggle the dark mode class
  body.classList.toggle("dark-mode");

  // Update the button text and save the mode in a cookie
  if (body.classList.contains("dark-mode")) {
    // button.innerHTML = "Light Mode";
    button.innerHTML = loadSvg("sun", 14);
    document.cookie = "darkMode=enabled; path=/; max-age=31536000"; // 1 year
  } else {
    // button.innerHTML = "Dark Mode";

    button.innerHTML = loadSvg("moon-star", 14);
    document.cookie = "darkMode=disabled; path=/; max-age=31536000"; // 1 year
  }
}

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
}

function applyDarkMode() {
  const body = document.body;
  const button = document.getElementById("dark-mode-button");

  // Check the cookie for the dark mode preference
  const darkMode = getCookie("darkMode");

  if (darkMode === "enabled") {
    body.classList.add("dark-mode");
    button.innerHTML = loadSvg("sun", 14);
  } else {
    body.classList.remove("dark-mode");
    button.innerHTML = loadSvg("moon-star", 14);
  }
}

// Call this function when the page loads
document.addEventListener("DOMContentLoaded", applyDarkMode);
