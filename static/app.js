new Vue({
    el: '#app',
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                //console.log(response.data[0].Priority)
                Vue.set(this, "tableData",
                    [{
                        date: '2016-05-03',
                        name: 'Tom',
                        address: 'No. 189, Grove St, Los Angeles'
                    }, {
                        date: '2016-05-02',
                        name: 'Tom',
                        address: 'No. 189, Grove St, Los Angeles'
                    }]
                )
            }).catch((error) => { console.log(error); });
        },
        handleEdit(index, row) {
            console.log(index, row);
        },
        handleDelete(index, row) {
            console.log(index, row);
        }
    },
    data: function() {
        return {
            tableData: [],
            visible: false
        }
    }
});