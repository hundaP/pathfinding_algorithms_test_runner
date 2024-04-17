import os
import glob
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from matplotlib.ticker import MaxNLocator

files = glob.glob('./data/averages*.csv')
all_data = []

for file in files:
    data = pd.read_csv(file)
    maze_size = int(os.path.basename(file).split('x')[0].replace('averages', ''))
    data['MazeSize'] = maze_size
    all_data.append(data)

all_data = pd.concat(all_data)
data_grouped = all_data.groupby(['Algorithm', 'MazeSize'])[['Time', 'MemoryUsed']].mean().reset_index()
data_dijkstra_astar = data_grouped[data_grouped['Algorithm'].isin(['dijkstra', 'astar'])]

data_other = data_grouped[~data_grouped['Algorithm'].isin(['dijkstra', 'astar'])]
colors = plt.get_cmap('tab10')
color_map = {algorithm: colors(i) for i, algorithm in enumerate(data['Algorithm'].unique())}

for data, group in [(data_dijkstra_astar, 'dijkstra_astar'), (data_other, 'other')]:
    for metric in ['Time', 'MemoryUsed']:
        fig, ax = plt.subplots(figsize=(15, 10))  # Increase the figure size
        sns.lineplot(data=data, x='MazeSize', y=metric, hue='Algorithm', ax=ax, palette=color_map)
        ax.set_title(f'{group} - {metric}')
        if metric == 'Time':
            if group == 'other':
                ax.set_ylim(bottom=0, top=400)
            ax.yaxis.set_major_locator(MaxNLocator(nbins=15))
        else:
            ax.set_ylim(bottom=0, top=data_grouped[metric].max() * 1.1)

        ax.spines['top'].set_visible(False)
        ax.spines['right'].set_visible(False)
        plt.subplots_adjust(left=0.05, right=0.95, top=0.95, bottom=0.05)
        plt.savefig(f'graphs/{group}_{metric}.svg', format='svg')