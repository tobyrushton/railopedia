<script lang="ts">
    import stations from '../../../../data/station-list.json'
    export let placeholder: string
    export let value: journey.IStation

    let filteredStations: journey.IStation[] = []

    const filterStations = (string: string): void => {
        filteredStations = stations.filter((station) => {
            return station.name.toLowerCase().includes(string.toLowerCase())
                || station.code.toLowerCase().includes(string.toLowerCase())
        }).slice(0, 5)
    }

    let inputValue = ""

    const handleSelect = (station: journey.IStation): void => {
        inputValue = station.name
        value = station
        filteredStations = []
    }

    $: if(inputValue === "") {
        filteredStations = []
    }

</script>

<div class="w-80">
    <input 
        class="border border-solid p-2 rounded placeholder:text-black w-full" 
        type="text"
        placeholder={placeholder}
        bind:value={inputValue}
        on:input={() => filterStations(inputValue)}
    />
    <div class="flex flex-col rounded-md shadow absolute bg-white z-10 w-80">
        {#if filteredStations.length > 0}
                {#each filteredStations as station}
                    <button 
                        class="flex gap-2 p-2 hover:bg-primary hover:text-white" 
                        on:click={() => handleSelect(station)}
                        tabindex="0"
                    >
                        <p class="text-gray-500">
                            {station.code}
                        </p>
                        <p>
                            {station.name}
                        </p>
                    </button>
                {/each}
        {/if}
    </div>
</div>