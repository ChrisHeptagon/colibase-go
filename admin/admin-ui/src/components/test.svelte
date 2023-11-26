<script lang="ts">
    let response = fetch('https://jsonplaceholder.typicode.com/todos/1').then(res => res.json());
    console.log(response);
</script>

<h1>Test</h1>
<div>
    {#await response}
        <p>loading...</p>
    {:then value}
        <p>{value.title}</p>
    {:catch error}
        <p>error {error.message}</p>
    {/await}
        
</div>