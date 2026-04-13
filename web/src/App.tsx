import "./App.css";

function App() {
  return (
    <div className="site-shell">
      <header className="hero" id="top">
        <p className="eyebrow">naterpatater.com</p>
        <h1>Nathan Brooks</h1>
      </header>

      <main className="content">
        <section className="section" aria-labelledby="projects-title">
          <h2 id="projects-title">Projects</h2>
          <article className="project-card">
            <h3>Checkboxes</h3>
            <p>
              ☑️ A lot of checkboxes inspired by
              https://onemillioncheckboxes.com/
            </p>
            <a href="https://checkboxes.naterpatater.com">Open Checkboxes</a>
          </article>
        </section>
      </main>

      <footer className="site-footer">
        <p>
          Favicon uses Twemoji via{" "}
          <a href="https://favicon.io/" target="_blank" rel="noreferrer">
            favicon.io
          </a>{" "}
          under{" "}
          <a
            href="https://creativecommons.org/licenses/by/4.0/"
            target="_blank"
            rel="noreferrer"
          >
            CC-BY 4.0
          </a>
          .
        </p>
      </footer>
    </div>
  );
}

export default App;
