from flask import Flask
from flask import render_template
import paramiko
import paramiko.util
from data import structured_data, aggregate_reports

app = Flask(__name__)

paramiko.util.log_to_file("paramiko.log")


@app.route('/', methods=['GET'])
def create():
    # data.generate_report_plot(data.report)
    aggregate_reports()
    dates = sorted(structured_data.keys())
    return render_template('index.html', dates=dates)


@app.route('/data/<date>')
def data_by_date(date):
    return structured_data.get(date)


if __name__ == '__main__':
    app.run(debug=True, port=5001)
