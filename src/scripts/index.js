//import axios from 'axios'

var day_buttons = new Vue({
    el:'#day_buttons',

    data: {
        days: [
            {
                id: 1,
                date: "2020-11-18T21:00:00.000Z",
                average_color: [255, 0, 0]
            },
            {
                id: 2,
                date: "2020-11-21T21:00:00.000Z",
                average_color: [255, 0, 0]
            },
            {
                id: 1,
                date: "2020-10-23T21:00:00.000Z",
                average_color: [255, 0, 0]
            },
            {
                id: 1,
                date: "2020-10-30T21:00:00.000Z",
                average_color: [255, 0, 0]
            },
            {
                id: 1,
                date: "2020-11-12T21:00:00.000Z",
                average_color: [255, 0, 0]
            },
            {
                id: 1,
                date: "2020-11-19T21:00:00.000Z",
                average_color: [255, 0, 0]
            }
        ]
    },

    methods: {
        colorToHex: function (color) {
            function compToHex(c) {
                let hex = c.toString(16);
                return hex.length === 1 ? "0" + hex : hex;
            }
            return "#" + compToHex(color[0]) + compToHex(color[1]) + compToHex(color[2]);
        },

        getDayUrl: function (day) {
            return "/" + day.getFullYear() + "_" + day.getMonth() + "_" + day.getDate();
        },

        isDayFilled: function (day) {
            for(let i = 0; i < this.days.length; ++i)
            {
                if(this.days[i].date === day.toISOString())
                    return true;
            }
            return false;
        },

        getFilledDayColor: function (day) {
            for(let i = 0; i < this.days.length; ++i)
            {
                if(this.days[i].date === day.toISOString())
                    return this.days[i].average_color;
            }
            return [169, 169, 169];
        },

        getDisplayedWeeks: function () {
            let closest_sunday = new Date();
            closest_sunday.setDate(closest_sunday.getDate() + 7 - closest_sunday.getDay());
            let three_weeks = [];
            for(let i = 2; i >= 0; --i)
            {
                let current_week = [];
                for(let j = 6; j >= 0; --j)
                {
                    current_week.push(new Date());
                    current_week[6 - j].setDate(closest_sunday.getDate() - j - i * 7);
                    current_week[6 - j].setHours(0, 0, 0, 0);
                }
                three_weeks.push(current_week);
            }
            return three_weeks;
        }
    }
    /*
    mounted() {
        axios.get("/api/days/brief").then(response => {
            this.days = response.data;
        })
    }*/

})