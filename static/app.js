new Vue({
    el: '#app',
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                Vue.set(this, "tableData", response.data)
            }).catch((error) => { console.log(error); });
        },
        handleEdit(index, row) {
            console.log(index, row);
        },
        handleDelete(index, row) {
            console.log(index, row);
            axios.delete("/rule/"+row.Id).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        }
    },
    data: function() {
        return {
            tableData: [],
            visible: false
        }
    },
    created: function() {
        this.gets()
    }
});