import "./app.css";
import { useState } from "preact/hooks";

export const App = () => {
  const [status, setStatus] = useState("idle");

  const handleSubmit = (e) => {
    setStatus("loading");
    e.preventDefault();

    const form = e.target;
    const formData = new FormData(form);

    const formJson = Object.fromEntries(formData.entries());
    console.log(formJson); // TODO: send
    fetch("/send", {
      method: "POST",
      body: JSON.stringify(formJson),
    }).then(
      () => {
        setStatus("completed");
      },
      (err) => {
        console.error(err);
        setStatus("error");
      }
    );
  };

  return (
    <Wrapper isPreview={location.hash.includes("preview")}>
      <form method="post" onSubmit={handleSubmit}>
        <label>Navn:</label>
        <input type="text" name="name" required />

        <div className="spacer" />

        <label>E-post:</label>
        <input type="email" name="email" required />

        <div className="spacer" />

        <label>Melding:</label>
        <textarea name="message" required />

        <div className="spacer" />

        <div class="button-wrapper">
          <button
            type="submit"
            className={status === "idle" ? "active enabled" : ""}
          >
            Send
          </button>
          <button
            type="submit"
            className={status === "error" ? "active enabled" : ""}
          >
            Prøv igjen
          </button>
          <button
            type="button"
            className={status === "loading" ? "active" : ""}
          >
            Sender...
          </button>
          <button
            type="button"
            className={status === "completed" ? "active" : ""}
          >
            Takk for din henvendelse!
          </button>
        </div>
      </form>
    </Wrapper>
  );
};

const Wrapper = ({ children, isPreview }) =>
  isPreview ? <div className="preview">{children}</div> : children;
