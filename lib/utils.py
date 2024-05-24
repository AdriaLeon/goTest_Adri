import matplotlib.colors as mcolors
from matplotlib import patches, pyplot as plt
import numpy as np
import json
import plotly.graph_objects as go

instrument_names = ['Voice', 'Violin', 'Mridangam', 'Ghatam']

def parse_json_file(file_name: str):
    file_path = "data/" + file_name
    # Open and read the JSON file
    with open(file_path, 'r') as file:
        data = json.load(file)
    
    # Extract the boolean interval arrays
    is_voice = data['is_voice']
    is_violin = data['is_violin']
    is_mridangam = data['is_mridangam']
    is_ghatam = data['is_ghatam']
    
    # Determine the number of intervals (assuming all arrays have the same length)
    num_intervals = len(is_voice)
    
    # Create a NumPy boolean matrix to store the intervals
    intervals_matrix = np.zeros((4, num_intervals), dtype=bool)
    
    # Populate the matrix with the intervals data
    intervals_matrix[0] = is_voice
    intervals_matrix[1] = is_violin
    intervals_matrix[2] = is_mridangam
    intervals_matrix[3] = is_ghatam
    
    return intervals_matrix

def plot_audio_from_json(json_file, save_name):
    with open(json_file, 'r') as file:
        data = json.load(file)

    audio_array = data['audio_array']

    x = list(range(len(audio_array)))
    y = audio_array

    plt.plot(x,y)
    plt.savefig(save_name, format='html')
    plt.close()

def save_boolean_matrix(bool_matrix, save_name="output.png"):
    n, m = bool_matrix.shape
    # Create a new matrix with spacing between the original matrix rows
    matrix_with_spacing = np.zeros((n*2, m), dtype=bool)
    matrix_with_spacing[::2] = bool_matrix
    
    # Define the colormap
    cmap = mcolors.ListedColormap(['white', 'darkgrey'])
    
    # Create the plot
    plt.imshow(matrix_with_spacing, cmap=cmap, interpolation='nearest', aspect='auto')
    
    # Set the x-axis and y-axis labels
    plt.xlabel('Time')
    plt.ylabel('Instruments')
    
    # Set the y-ticks to the instrument names
    plt.yticks(np.arange(0, 2 * n, 2), instrument_names)
    
    # Set the title of the plot
    plt.title('Result of Instrument Classification')
    
    # Save the plot to a file
    plt.savefig("public/" + save_name)
    plt.close()

def save_boolean_matrix_interactive(bool_matrix, save_name="output.html"):
    n, m = bool_matrix.shape
    # Create a new matrix with spacing between the original matrix rows
    matrix_with_spacing = np.zeros((n*2, m), dtype=bool)
    matrix_with_spacing[::2] = bool_matrix
    
    # Define the colormap
    cmap = ['white', 'darkgrey']
    
    # Create the figure
    fig = go.Figure()

    # Add rectangles for each True value in the matrix
    for i in range(n):
        for j in range(m):
            if matrix_with_spacing[i * 2, j]:
                fig.add_shape(
                    type="rect",
                    x0=j,
                    y0=i * 2,
                    x1=j + 1,
                    y1=i * 2 + 1,
                    fillcolor=cmap[1],
                    line=dict(width=0),
                )
    
    # Set the layout
    fig.update_layout(
        title='Result of Instrument Classification',
        xaxis=dict(title='Timestamps'),
        yaxis=dict(title='Instruments', tickmode='array', tickvals=list(range(0, 2 * n, 2)), ticktext=instrument_names),
        yaxis_autorange='reversed',
        showlegend=False
    )

    # Save the plot to an HTML file
    fig.write_html("public/" + save_name)

def generate_waveform(file_name, save_name="output.html"):
    # Load JSON data
    file_path = "data/" + file_name
    # Open and read the JSON file
    with open(file_path, 'r') as file:
        data = json.load(file)

    # Extract amplitude values
    amplitudes = data["audio_array"]

    # Create a time axis based on the number of amplitude samples
    time = np.linspace(0, len(amplitudes), num=len(amplitudes))

    # Plot the sound wave
    plt.figure(figsize=(10, 4))
    plt.plot(time, amplitudes, label='Sound Wave')
    plt.xlabel('Time')
    plt.ylabel('Amplitude')
    plt.title('Sound Wave Visualization')
    plt.legend()
    plt.grid(True)
    plt.savefig('public/sound_wave.png')  # Save the plot as an image file
    plt.close()  # Close the plot to free up memory