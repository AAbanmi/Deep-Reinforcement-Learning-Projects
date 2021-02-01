# Description of implementation:
### Files and their usages:
There are 3 main files needed for the agent to train:
1. Navigation.ipynb (Contains: Running the environment, agent, the deep Q-learning and plots showing the results)
2. model.py (Contains: The function and implementation of the Q network)
3. dqn_agent.py (Contains: The agent implememtations used in Navigation.ipynb)


### Learning algorithm:
We are usimg Q-Learning algorithm, which is a reinforcement learning algorithm with a Q-table saving the action/state combination.

### Hyperparameters:
1. n_episodes (int): maximum number of training episodes
2. max_t (int): maximum number of timesteps per episode
3. eps_start (float): starting value of epsilon, for epsilon-greedy action selection
4. eps_end (float): minimum value of epsilon
5. eps_decay (float): multiplicative factor (per episode) for decreasing epsilon
6. gamma (float): discount factor
7. batch_size (int): size of each training batch

### Architecture of the neural network:
We have 3 layers:
1. Input layer with the same size as state_size (37)
2. Hidden layer with size of 64
3. Output layer with the same size as action_size (4)

All of these layers are using activation function ReLu.


### Plot of the results:
Episode 100	Average Score: 4.92	Highest Score: 16.0
Episode 200	Average Score: 8.63	Highest Score: 20.0
Episode 300	Average Score: 12.48	Highest Score: 21.0
Episode 314	Average Score: 13.03	Highest Score: 24.0
Environment solved in 314 episodes!	Average Score: 13.03	Highest Score: 24.0

Plot of the results can be found in the image named "Project-1-Results.png" in this folder.

### Future improvement ideas:
Some of the ideas that could enhance the agents performance are using the following:

1. Experience replay
2. Double DQN
3. Dueling DQN
4. Prioritized Experience Replay
5. Rainbow approach
 
