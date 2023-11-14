<script lang="ts">
    import LoginForm from "$lib/components/loginForm.svelte";
  import { page } from "$app/stores";
</script>

<svelte:head>
    <title>Colibase Admin Login</title>
    
</svelte:head>

<LoginForm Heading="Login" UserSchema="{$page.data.UserSchema}"/>

