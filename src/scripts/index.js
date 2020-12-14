var dayButtons = new Vue({
    el: "#dayButtons",

    data: {
        days: []
    },

    methods: {
        RGBColorToHex(color) {
            function convertToHex(c) {
                let hex = c.toString(16);
                return hex.length === 1 ? "0" + hex : hex;
            }
            return "#" + convertToHex(color[0]) + convertToHex(color[1]) + convertToHex(color[2]);
        },

        isDayFilled(day) {
            for(let filledDay of this.days) {
                if(Date.parse(filledDay.date) === Date.parse(day.toISOString())) {
                    return true;
                }
            }
            return false;
        },

        getFilledDayColor(day) {
            for(let filledDay of this.days) {
                if(Date.parse(filledDay.date) === Date.parse(day.toISOString())) {
                    return this.RGBColorToHex(filledDay.average_color);
                }
            }
            return "#a9a9a9";
        },

        getDisplayedWeeks() {
            let closestSunday = new Date();
            closestSunday = new Date(Date.UTC(closestSunday.getFullYear(), closestSunday.getMonth(),
                closestSunday.getDate() + 7 - closestSunday.getDay()));
            let threeWeeks = Array.from(Array(3), () => new Array(7));
            for(let i = 0; i < 3; ++i) {
                for(let j = 0; j < 7; ++j) {
                    threeWeeks[i][j] = new Date(closestSunday);
                    threeWeeks[i][j].setUTCDate(closestSunday.getUTCDate() - 20 + i * 7 + j);
                }
            }
            return threeWeeks;
        }
    },

    mounted() {
        axios.get("/api/days/brief").then((response) => {
            this.days = response.data;
        });
    }

});
