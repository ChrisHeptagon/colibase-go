<script lang="ts">
export const ssr = true
import ColibaseLogo from "$lib/assets/favicon.ico?raw";
</script>

<slot />



<svelte:head>
  <link rel="favicon" href={ColibaseLogo} />
</svelte:head>
