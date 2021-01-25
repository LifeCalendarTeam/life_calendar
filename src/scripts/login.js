var login = new Vue({
    el:'#loginform',

    data: {
        user_id: "",
        password: ""
    },

    methods: {
        trySignIn() {
            if(!isNaN(this.user_id) && this.user_id !== "" && this.password !== "") {
                axios({
                    method: 'post',
                    url: '/login',
                    params: {
                        user_id: parseInt(this.user_id, 10),
                        password: this.password
                    },
                    headers: {
                        "Content-type": "application/x-www-form-urlencoded"
                    }
                })
                .then((response) => {
                    console.log(response);
                    document.location.href = "/";
                })
                .catch(function(error) {
                    console.log(error);
                    alert("An error occurred! Please try again or contact the authors.");
                });
            }
            else if(this.user_id === "" || this.password === "") {
                alert("Both username and password are required. Please try again.");
            }
            else if(!isNaN(this.user_id)) {
                console.error(error);
            }
        }
    }
})