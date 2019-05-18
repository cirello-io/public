const reducer = (state = {}, action) => {
  switch (action.type) {
    case 'ADD_LINK':
      return { ...state, loading: true }
    case 'LINK_ADDED':
      console.log("link added", action.link)
      return { ...state, addedLink: action.link, loading: false }
    default:
      return state
  }
}
export default reducer
