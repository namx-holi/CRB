
import subprocess
import time
import matplotlib.pyplot as plt
import numpy as np


def estimate_completion_times(counts, loops, cooldowns, resp_times):
	print("--Estimated completion times--")
	for resp_time in resp_times:
		estimated_time = loops * (cooldown + resp_time*1000) * len(counts) - cooldown

		print("  %4dms average response time will take %d seconds (%.2f minutes)" % (
				resp_time*1000, estimated_time/1000, estimated_time/60000))
	print("\n")


def run_full_benchmark(url, counts, loops=5, cooldown=3000):
	estimate_completion_times(counts, loops, cooldown, [0, 1, 2, 3])

	results = {}

	for i, count in enumerate(counts):
		print("Running for {} requests (Benchmark {} of {})".format(
				count, i+1, len(counts)))
		output = run_crb(url, loops=loops, count=count, cooldown=cooldown)

		# sleep for a cooldown. this is because there's no cooldown time
		# at the end of a run of crb
		if i != len(counts)-1:
			time.sleep(cooldown/1000.0) #divided by 1000 because sleep takes seconds

		results[count] = get_results_from_output(output)

	print("Done!")

	return results


def run_crb(url, count=1, loops=1, cooldown=1000, verbose=False, show_command=False):
	command = ["./crb"]

	if url:
		command.append("-url=%s" % url)
	
	if count:
		command.append("-count=%d" % count)
	
	if loops:
		command.append("-loops=%d" % loops)
	
	if cooldown:
		command.append("-cooldown=%d" % cooldown)
	
	if verbose:
		command.append("-verbose")

	command.append("-display")

	if show_command:
		print("We are running:\n  %s\n" % " ".join(command))

	p = subprocess.run(command, stdout=subprocess.PIPE)

	return p.stdout.decode("utf-8")


def get_results_from_output(output):
	# separeate into lines
	lines = output.split("\n")

	# find from where overall stats are
	header_index = lines.index("--DISPLAYING RESULTS--")

	# empty results dict
	results_as_string = lines[header_index + 1].split(", ")

	return [float(result) for result in results_as_string]


# UNUSED
def results_to_csv(stats, counts, looking_for):
	# counts are added so data can be in order

	delim = ","

	rows = []

	# header
	header_items = counts
	header_items.reverse()

	"""
	header will look like

	1,2,3,4,5,10,15,20,25,....

	will look like data but is the number of requests
	"""
	rows.append(delim.join(header_items))

	for count in counts:
		line_items = [count]

		for key in looking_for:
			line_items.append(stats[count][key])

		rows.append(delim.join(
			str(line_item) for line_item in line_items))

	return "\n".join(rows)


# UNUSED
def write_csv_to_file(csv_text, filename):
	f = open(filename, "w")
	f.write(csv_text)
	f.close()


def boxplot_results(results, counts):
	data = []

	for count in counts:
		x = np.array(results[count])

		# keep only the good points
		# ~ operates as logical not operator on boolean
		# numpy arrays
		filtered = x[~is_outlier(x, 1.5)]

		data.append(filtered)

	boxplot_scale_factor = 0.75
	if len(counts) > 1:
		widths = []

		# for each count
		for i in range(len(counts)):
			# distances between this count and adjacent counts
			left_between  = counts[(i) % len(counts)] - counts[(i-1) % len(counts)]
			right_between = counts[(i+1) % len(counts)] - counts[(i) % len(counts)]

			widths.append(min(abs(left_between), abs(right_between))*boxplot_scale_factor)
	else:
		widths = [1]

	f, ax = plt.subplots()
	ax.boxplot(data, positions=counts, widths=widths,
			# patch_artist=True,
			# showmeans=True,
			meanline=True)

	# make sure we start at 0 for each axis please
	# for xmax, find the last count and add the last width on so there's
	# space for the last boxplot
	xmax = counts[-1] + widths[-1]
	ax.set_xlim(xmin=0, xmax=xmax)
	ax.set_ylim(ymin=0)

	# set titles
	ax.set_title("Response times for concurrent requests")
	ax.set_xlabel("No. concurrent requests")
	ax.set_ylabel("Response time (ms)")

	plt.show()


#?????????
def is_outlier(points, thresh=3.5):
    """
    Returns a boolean array with True if points are outliers and False 
    otherwise.

    Parameters:
    -----------
        points : An numobservations by numdimensions array of observations
        thresh : The modified z-score to use as a threshold. Observations with
            a modified z-score (based on the median absolute deviation) greater
            than this value will be classified as outliers.

    Returns:
    --------
        mask : A numobservations-length boolean array.

    References:
    ----------
        Boris Iglewicz and David Hoaglin (1993), "Volume 16: How to Detect and
        Handle Outliers", The ASQC Basic References in Quality Control:
        Statistical Techniques, Edward F. Mykytka, Ph.D., Editor. 
    """
    if len(points.shape) == 1:
        points = points[:,None]
    median = np.median(points, axis=0)
    diff = np.sum((points - median)**2, axis=-1)
    diff = np.sqrt(diff)
    med_abs_deviation = np.median(diff)

    modified_z_score = 0.6745 * diff / med_abs_deviation

    return modified_z_score > thresh





if __name__ == "__main__":
	print("Imports complete\n\n")

	generate_count = False

	# settings for intervals if generating
	min_value = 0
	max_value = 10
	subvalues = 5

	if generate_count:
		counts = [
			i*int((max_value-min_value)/subvalues) + min_value
			for i in range(1, 1+subvalues)]
	
	else:
		counts = [1, 2, 3, 4, 5, 10, 15, 20, 30, 40, 50, 75, 100, 150, 200]
		# counts = [20]

	print("Counts:", counts, "\n\n")




	url = "http://localhost"
	
	loops = 5
	cooldown = 500

	results = run_full_benchmark(
			url, counts,
			loops=loops, cooldown=cooldown)

	boxplot_results(results, counts)
