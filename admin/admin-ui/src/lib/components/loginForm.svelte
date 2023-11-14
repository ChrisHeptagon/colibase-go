
<script lang="ts">
  import ColibaseLogo from "$lib/assets/colibase_logo.svg"
  import { enhance } from "$app/forms";
  export let UserSchema: any;
  let error: string | Record<string, unknown> | undefined;
  export let Heading: string;
  let showError: boolean = false;

  const handleInputChange = (e: Event) => {
    const input = e.target as HTMLInputElement;
    const errorSpan = document.querySelector(
      `#${input.name}-error`
    ) as HTMLSpanElement;
    if (input.validity.valid === false) {
      if (errorSpan) {        
      if (input.validity.valueMissing) {
        errorSpan.innerText = `Please enter a ${input.name}`;
      } else if (input.validity.patternMismatch) {
        errorSpan.innerText = `Please enter a valid ${input.name}`;
      } else if (input.validity.tooShort) {
        errorSpan.innerText = `Please enter a longer ${input.name}`;
      } else if (input.validity.tooLong) {
        errorSpan.innerText = `Please enter a shorter ${input.name}`;
      } else {
        errorSpan.innerText = `Please enter a valid ${input.name}`;
      }
    }
    else {
      console.log("error span not found");
    }
    } else {
      errorSpan.innerText = "";
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    const form = e.target as HTMLFormElement;
    form.querySelectorAll("input").forEach((input) => {
      handleInputChange({ target: input } as unknown as Event);
    });
  };
</script>

<div class="login-form">
<form id="login-form" novalidate autocomplete="on" method="POST" action="?/action" 	use:enhance={({ formElement, formData, action, cancel, submitter }) => {
  formElement.addEventListener("submit", handleSubmit);
  return async ({ result, update }) => {
    console.log("result: ", result);
    if (result) {
      if (result.type === "failure") {
        error = result.data
        console.log("error: ", error);
      } else if (result.type === "success") {
        error = "";
      } else if (result.type === "error") {
        error = result.error;
      }
    }
  }
}}>
  <h1>{Heading}</h1>
  <img class="logo" src="{ColibaseLogo}" alt="logo" />
  <fieldset>
    <legend>Enter your credentials</legend>
    {#if UserSchema}
      {#each UserSchema as field}
        <label for={field.name}>{field.name}</label>
        <input
          type={field.name}
          name={field.name}
          id={field.name}
          required={field.required}
          autocomplete="on"
          aria-required={field.required}
          aria-label={field.name}
          aria-describedby={`${field.name}-error`}
          aria-invalid={showError}
        />
        <span id={`${field.name}-error`} class="form-error"></span>
      {/each}
    {/if}
      <p class="error-item">
        {#if error}
          {error}
        {/if}
      </p>
  </fieldset>
  <button type="submit" id="button">Login</button>
</form>
</div>

<style lang="scss">
  @keyframes fadeIn {
    0% {
      opacity: 0;
    }
    100% {
      opacity: 1;
    }
  }
  @keyframes fadeOut {
    0% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
  * {
    font-family: Arial, Helvetica, sans-serif;
  }
  .login-form {
    animation: fadeIn 0.5s ease-in-out;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.75);
  }
  form {
    background-color: #ffffff;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.75);
  }
  .logo {
    margin: 0 auto;
    margin-bottom: 20px;
    display: block;
    width: 100px;
    filter: drop-shadow(0px 0px 10px rgba(0, 0, 0, 0.75));
  }
  h1 {
    text-align: center;
    color: #000000;
    font-size: 24px;
    margin: 0 0 20px 0;
    text-shadow: 0px 0px 10px rgba(0, 0, 0, 0.5);
  }
  fieldset {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: none;
    padding-left: 10px;
    padding-right: 10px;
  }
  legend {
    font-size: 20px;
    text-align: center;
  }

  button {
    border: none;
    border-radius: 5px;
    padding: 10px;
    font-weight: bold;
    width: 100%;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.75);
    font-size: 16px;
    background-image: radial-gradient(
      circle at 90% 100%,
      #b8e986 0%,
      #8ecf70 50%,
      #5ca860
    );
    cursor: pointer;
  }
  @keyframes outlineFadeIn {
    0% {
      outline: 0px solid #000000;
    }
    100% {
      outline: 2px solid #000000;
    }
  }
  @keyframes outlineFadeOut {
    0% {
      outline: 2px solid #000000;
    }
    100% {
      outline: 0px solid #000000;
    }
  }
  @keyframes activeButton {
    0% {
      filter: brightness(1);
    }
    100% {
      filter: brightness(0.8);
    }
  }
  button:hover {
    outline: 2px solid #000000;
    animation: outlineFadeIn 0.5s ease-in-out;
  }
  button:not(:hover) {
    animation: outlineFadeOut 0.5s ease-in-out;
    outline: none;
  }
  button:active {
    animation: activeButton 0.5s ease-in-out;
    filter: brightness(0.8);
  }
  button:not(:active) {
    transition: filter 0.5s ease-in-out;
    filter: brightness(1);
  }
  .error-item:not(:empty) {
    color: rgb(0, 0, 0);
    text-transform: capitalize;
    margin: 10px;
    height: 20px;
    border-color: red;
    border-style: solid;
    border-radius: 5px;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.75);
    padding: 10px;
    font-size: 14px;
    transition: all 0.5s ease-in-out;
  }
  .error-item:empty {
    transition: all 0.5s ease-in-out;
    color: transparent;
    margin: 0;
    padding: 0;
    height: 0;
    width: 0;
  }


  label {
    font-size: 16px;
    margin-top: 10px;
    align-self: flex-start;
    text-shadow: 0px 0px 10px rgb(0, 0, 0, 0.5);
    display: block;
  }
  input {
    display: block;
    box-shadow: 0px 0px 10px 0px rgb(0, 0, 0, 0.5);
    border: 1px 2px solid #000000;
    border-radius: 5px;
    padding: 10px;
    margin-top: 5px;
    margin-bottom: 5px;
    font-size: 16px;
  }

  .form-error:not(:empty) {
    color: red;
    font-size: 14px;
    text-shadow: 0px 0px 10px rgb(92, 11, 11);
    margin: 5px;
    align-self: flex-start;
    transition: all 0.5s ease-in-out;
  }

  .form-error:empty {
    transition: all 0.5s ease-in-out;
    color: transparent;
    margin: 0;
    padding: 0;
    height: 0;
    width: 0;
  }

</style>
