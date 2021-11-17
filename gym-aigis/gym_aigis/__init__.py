from gym.envs.registration import register

register(
    id='aigis-v0',
    entry_point='gym_aigis.envs:AigisEnv',
)