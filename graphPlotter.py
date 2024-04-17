import os
import glob
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from matplotlib.ticker import MaxNLocator
import numpy as np

# Get a list of all CSV files
files = glob.glob('averages*.csv')

# Initialize an empty list to hold all the data
all_data = []

# Loop over the files
for file in files:
    # Load the data from the file
    data = pd.read_csv(file)
    
    # Extract the maze size from the file name
    maze_size = int(os.path.basename(file).split('x')[0].replace('averages', ''))
    
    # Add a 'MazeSize' column to the data
    data['MazeSize'] = maze_size
    
    # Append the data to the all_data list
    all_data.append(data)

# Concatenate all the data into a single DataFrame
all_data = pd.concat(all_data)

# Group the data by 'Algorithm' and 'MazeSize' and calculate the mean 'Time' and 'MemoryUsed'
data_grouped = all_data.groupby(['Algorithm', 'MazeSize'])[['Time', 'MemoryUsed']].mean().reset_index()

# Separate the data for the 'dijkstra' and 'astar' algorithms
data_dijkstra_astar = data_grouped[data_grouped['Algorithm'].isin(['dijkstra', 'astar'])]

# Separate the data for the other algorithms
data_other = data_grouped[~data_grouped['Algorithm'].isin(['dijkstra', 'astar'])]
colors = plt.get_cmap('tab10')
color_map = {algorithm: colors(i) for i, algorithm in enumerate(data['Algorithm'].unique())}

# Create separate graphs for each group of algorithms
for data, group in [(data_dijkstra_astar, 'dijkstra_astar'), (data_other, 'other')]:
    for metric in ['Time', 'MemoryUsed']:
        # Plot the actual data
        fig, ax = plt.subplots(figsize=(15, 10))  # Increase the figure size
        sns.lineplot(data=data, x='MazeSize', y=metric, hue='Algorithm', ax=ax, palette=color_map)
        ax.set_title(f'{group} - {metric} (actual)')
        if metric == 'Time':
            if group == 'other':
                ax.set_ylim(bottom=0, top=350)
            ax.yaxis.set_major_locator(MaxNLocator(nbins=15))  # Set the maximum number of ticks on the y-axis
        else:
            ax.set_ylim(bottom=0, top=data_grouped[metric].max() * 1.1)  # Set the lower and upper limits of y-axis

        # Remove the top and right spines
        ax.spines['top'].set_visible(False)
        ax.spines['right'].set_visible(False)

        # Reduce the borders further
        plt.subplots_adjust(left=0.05, right=0.95, top=0.95, bottom=0.05)

        # Save the figure as an SVG file
        plt.savefig(f'graphs/{group}_{metric}_actual.svg', format='svg')