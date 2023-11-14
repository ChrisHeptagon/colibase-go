<script lang="ts">
    import LoginForm from "$lib/components/loginForm.svelte";
	import type { PageData } from './$types';
	
	export let data: PageData;
    </script>



<LoginForm Heading="Initial Login" UserSchema="{data.UserSchema}"/>

