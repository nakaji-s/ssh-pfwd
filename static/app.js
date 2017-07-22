new Vue({
    el: '#app',
    data: {
    },
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                console.log(response.data[0].Priority)
            }).catch((error) => { console.log(error); });
        }
    },
    data: function() {
        return { visible: false }
    }
});